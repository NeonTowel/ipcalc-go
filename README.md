# ipcalc-go

A Go implementation of the popular `ipcalc` tool for calculating IP network information.

This is a port of the original Perl [ipcalc](https://github.com/kjokjo/ipcalc) tool, rewritten in Go.

## Features

- IPv4 and IPv6 support
- Network, broadcast, and host range calculation
- Subnet splitting
- IP range deaggregation
- Binary representation of addresses
- Colorized output
- HTML output option

## Installation

### Using Task

This project uses [Task](https://taskfile.dev) for build automation. To install Task, follow the instructions on the [Task website](https://taskfile.dev/#/installation).

Once Task is installed, you can use the following commands:

```bash
# List all available tasks
task -l

# Build only
task build

# Run tests
task test

# Install to your system (platform-specific)
task install

# Run tests, build, and install
task all

# Uninstall
task uninstall

# Clean build artifacts
task clean
```

The installation process is platform-specific:
- On Windows, it installs to `%USERPROFILE%\bin` and adds it to your PATH
- On macOS, it installs to `/usr/local/bin`
- On Linux, it installs to `/usr/local/bin` (requires sudo)

### From Source (Manual)

```bash
git clone https://github.com/neontowel/ipcalc-go.git
cd ipcalc-go
go build -o ipcalc ./cmd/ipcalc
# Copy the binary to a location in your PATH
```

## Usage

```
Usage: ipcalc [options] <ADDRESS>[[/]<NETMASK>] [NETMASK]

ipcalc takes an IP address and netmask and calculates the resulting
broadcast, network, Cisco wildcard mask, and host range. By giving a
second netmask, you can design sub- and supernetworks. It is also
intended to be a teaching tool and presents the results as easy-to-
understand binary values.

Options:
  -h, --help        Display help usage
  -n, --nocolor     Don't display ANSI color codes
  -b, --nobinary    Suppress the bitwise output
  -c, --class       Just print bit-count-mask of given address
  -H, --html        Display results as HTML
  -v, --version     Print Version
  -s, --split       Split into networks of specified sizes
  -r, --range       Deaggregate address range
```

## Examples

### Basic IPv4 calculation

```bash
ipcalc 192.168.1.1/24
```

Output:
```
Address:   192.168.1.1          11000000.10101000.00000001.00000001
Netmask:   255.255.255.0 = 24   11111111.11111111.11111111.00000000
Wildcard:  0.0.0.255            00000000.00000000.00000000.11111111
=>
Network:   192.168.1.0/24       11000000.10101000.00000001.00000000
HostMin:   192.168.1.1          11000000.10101000.00000001.00000001
HostMax:   192.168.1.254        11000000.10101000.00000001.11111110
Broadcast: 192.168.1.255        11000000.10101000.00000001.11111111
Hosts/Net: 254                   Class C, Private Internet
```

### Basic IPv6 calculation

```bash
ipcalc fde6:36fc:c985:0:c2c1:c0ff:fe1d:cc7f 64
```

Output:
```
Address: fde6:36fc:c985::c2c1:c0ff:fe1d:cc7f     1111110111100110:0011011011111100:1100100110000101:0000000000000000:1100001011000001:1100000011111111:1111111000011101:1100110001111111
Netmask: 64                                      1111111111111111:1111111111111111:1111111111111111:1111111111111111:0000000000000000:0000000000000000:0000000000000000:0000000000000000
Prefix:  fde6:36fc:c985::/64                     1111110111100110:0011011011111100:1100100110000101:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000
```

### Deaggregating an IP range

```bash
ipcalc -r 192.168.0.1 192.168.0.10
```

### Splitting a network into subnets

```bash
ipcalc -s 192.168.0.0/24 100 50 25
```

## License

This project is licensed under the GPL License - see the LICENSE file for details.

## Acknowledgments

* Original ipcalc by Krischan Jodies
* Contributors to the original ipcalc project 