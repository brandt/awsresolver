# awsresolver

Very basic DNS server that resolves AWS's EC2 internal FQDNs (ex: `ip-192-0-2-1.us-west-2.compute.internal`) by extracting the IP out of the hostname.

Includes an optional subcommand that sets up macOS to automatically direct all "*.internal" queries to this daemon.


## Requirements

- Go 1.13 or later to build.

NOTE: Only macOS is currently supported.


## Installation

1. Run: `brew install brandt/personal/awsresolver`
2. Run: `sudo awsresolver setup` (installs `/etc/resolver/internal`)
3. Run: `brew services start awsresolver`

### Confirm it's working

This hooks into macOS's resolver. That means things like `ping` and `ssh` will do what you expect, but `dig` will not.

To confirm it's working, run: `ping ip-192-0-2-1.us-west-2.compute.internal`

If `awsresolver` is correctly setup, you will see ping attempt to reach `192.0.2.1`:

```
# SUCCESS                                     vvvvvvvvv
PING ip-192-0-2-1.us-west-2.compute.internal (192.0.2.1): 56 data bytes
Request timeout for icmp_seq 0
Request timeout for icmp_seq 1
```


## How it works

This tool listens for `A` record requests ending in `.internal`, extracts the IP from the requested name, and returns it as a response. It binds to UDP and TCP `127.0.0.1:1053`.

Mac OS X has a cool feature that allows you to configure different resolvers by domain. (See: `man 5 resolver`) When you run `sudo awsresolver setup`, it writes a config file to `/etc/resolver/internal` that steers `*.internal` requests to `127.0.0.1:1053`.

Note that the resolver(5) config only applies to DNS resolution performed through the built-in OS facilities.  So `ping`, `ssh`, and Chrome will be routed to this resolver, but by default `dig` will not.

To query with dig, point it directly at the resolver like so: `dig @127.0.0.1 -p 1053 ip-192-0-2-1.us-west-2.compute.internal`


## Building

To build from source, simply run these commands from inside this repo:

    go mod vendor # optional
    make

The compiled binary is here: `bin/awsresolver`


## Uninstalling

To uninstall:

1. Run: `sudo rm -f /etc/resolver/internal`
2. Run: `brew uninstall awsresolver`


## Author

- J. Brandt Buckley
