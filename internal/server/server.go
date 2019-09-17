package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/miekg/dns"
)

// ec2InternalHostnameMatcher matches the two variations of internal hostnames used
// by AWS.
//
// - US-EAST1 format:      ip-private-ipv4-address.ec2.internal
// - Other regions format: ip-private-ipv4-address.region.compute.internal
var ec2InternalHostnameMatcher = regexp.MustCompile(`(?i)^ip-(\d{1,3})-(\d{1,3})-(\d{1,3})-(\d{1,3})[.](?:[a-zA-Z0-9-]{1,255}[.]compute|ec2)[.]internal[.]?$`)

// Start the DNS server TCP and UDP listeners.
func Start() {
	dns.HandleFunc(".", handleRequest)

	// UDP listener
	go func() {
		srv := &dns.Server{Addr: "127.0.0.1:1053", Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to set udp listener %s", err.Error())
		}
	}()

	// TCP listener
	go func() {
		srv := &dns.Server{Addr: "127.0.0.1:1053", Net: "tcp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to set tcp listener %s", err.Error())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%d) received, stopping...\n", s)
}

// parseInternalDomainQuery takes a dns.Question and returns the IP if
// the given domain was an internal EC2 domain. If no such IP was found
// or the query type is incorrect, nil is returned.
func parseInternalDomainQuery(q dns.Question) net.IP {
	if q.Qclass != dns.ClassINET {
		return nil
	}

	switch q.Qtype {
	case dns.TypeA:
		matches := ec2InternalHostnameMatcher.FindStringSubmatch(q.Name)
		if len(matches) != 5 {
			return nil
		}
		return net.ParseIP(fmt.Sprintf("%s.%s.%s.%s", matches[1], matches[2], matches[3], matches[4]))
	case dns.TypeAAAA:
		// AWS is not currently setting hostnames on the IPv6 addresses:
		// "We do not provide DNS hostnames for IPv6 addresses."
		// https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html
		return nil
	default:
		return nil
	}
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// We set a relatively short TTL even though the encoded IPs will never
	// change in case we happen to catch a non-EC2 internal address. Also,
	// it's not like this is hard for us to compute.
	ttl := uint32(0)

	q := r.Question[0]

	log.Printf("QUERY: %s", q.String())

	m := new(dns.Msg)

	// AWS only sets IPv4 hostnames currently.
	ip := parseInternalDomainQuery(q)

	// If we can't get an IP out of it or it's the wrong query type, reply
	// with a SERVFAIL.
	// TODO: Is there a more appropriate response code?
	if ip == nil {
		m = m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}

	domain := q.Name
	m.SetReply(r)
	m.Authoritative = true

	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
	rr.A = ip

	m.Answer = []dns.RR{rr}
	w.WriteMsg(m)
}
