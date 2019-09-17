# awsresolver

Resolves AWS internal IPs like `ip-10-78-32-168.us-east-2.compute.internal`.

## Requirements

- Go 1.13 or later.

NOTE: Only macOS is currently supported.

## Installation

1. Install `go get -u github.com/brandt/awsresolver`
2. Setup daemon and `/etc/resolver/internal` config: `sudo awsresolver setup`

### Confirming it's working

This hooks into macOS's resolver. That means things like `ping` and `ssh` will do what you expect, but `dig` will not.

To confirm it's working, run: `ping ip-192-0-2-1.us-west-2.compute.internal`

If `awsresolver` is correctly setup, you will see ping attempt to reach `192.0.2.1`:

```
# SUCCESS
PING ip-192-0-2-1.us-west-2.compute.internal (192.0.2.1): 56 data bytes
Request timeout for icmp_seq 0
Request timeout for icmp_seq 1
```

## Author

- J. Brandt Buckley
