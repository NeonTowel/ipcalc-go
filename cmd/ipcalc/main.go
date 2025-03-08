package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neontowel/ipcalc-go/pkg/calculator"
	"github.com/neontowel/ipcalc-go/pkg/formatter"
	"github.com/spf13/pflag"
)

const version = "0.1.0"

func main() {
	// Define command-line flags
	help := pflag.BoolP("help", "h", false, "Display help usage")
	noColor := pflag.BoolP("nocolor", "n", false, "Don't display ANSI color codes")
	noBinary := pflag.BoolP("nobinary", "b", false, "Suppress the bitwise output")
	classOnly := pflag.BoolP("class", "c", false, "Just print bit-count-mask of given address")
	html := pflag.BoolP("html", "H", false, "Display results as HTML")
	showVersion := pflag.BoolP("version", "v", false, "Print Version")
	split := pflag.BoolP("split", "s", false, "Split into networks of specified sizes")
	deaggregate := pflag.BoolP("range", "r", false, "Deaggregate address range")

	// Parse flags
	pflag.Parse()

	// Get remaining arguments
	args := pflag.Args()

	// Check for help flag
	if *help || len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	// Check for version flag
	if *showVersion {
		fmt.Printf("ipcalc-go version %s\n", version)
		os.Exit(0)
	}

	// Set up output format
	format := formatter.OutputFormat{
		UseColor:  !*noColor && !*html && isTerminal(),
		UseHTML:   *html,
		UseBinary: !*noBinary,
	}

	// Print HTML header if needed
	if format.UseHTML {
		fmt.Print(formatter.FormatHTMLHeader())
		fmt.Printf("<!-- Version %s -->\n", version)
	}

	// Handle class-only mode
	if *classOnly && len(args) > 0 {
		handleClassOnly(args[0])
		os.Exit(0)
	}

	// Handle deaggregate mode
	if *deaggregate {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: Deaggregate mode requires two IP addresses")
			os.Exit(1)
		}
		handleDeaggregate(args[0], args[1], format)
		os.Exit(0)
	}

	// Handle split mode
	if *split {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: Split mode requires an IP/netmask and at least one size")
			os.Exit(1)
		}
		handleSplit(args[0], args[1:], format)
		os.Exit(0)
	}

	// Handle normal mode
	if len(args) > 0 {
		handleNormal(args, format)
	}

	// Print HTML footer if needed
	if format.UseHTML {
		fmt.Print(formatter.FormatHTMLFooter())
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println(`Usage: ipcalc [options] <ADDRESS>[[/]<NETMASK>] [NETMASK]

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

Examples:
  ipcalc 192.168.0.1/24
  ipcalc 192.168.0.1/255.255.128.0
  ipcalc 192.168.0.1 255.255.128.0 255.255.192.0
  ipcalc 192.168.0.1 0.0.63.255
  ipcalc -r 192.168.0.1 192.168.0.10
  ipcalc -s 192.168.0.0/24 10 20 30`)
}

// handleClassOnly handles the class-only mode
func handleClassOnly(ipStr string) {
	// Check if it's an IPv6 address
	if strings.Contains(ipStr, ":") {
		fmt.Println("IPv6 addresses don't have classes")
		return
	}

	// Parse the IP address
	ip, err := calculator.ParseIPv4(ipStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get the class
	class := calculator.GetClass(ip)
	bits := calculator.GetClassBits(class)

	// Print the result
	fmt.Println(bits)
}

// handleDeaggregate handles the deaggregate mode
func handleDeaggregate(startStr, endStr string, format formatter.OutputFormat) {
	// Check if these are IPv6 addresses
	if strings.Contains(startStr, ":") || strings.Contains(endStr, ":") {
		fmt.Fprintln(os.Stderr, "Error: IPv6 deaggregation is not supported yet")
		os.Exit(1)
	}

	// Deaggregate the range
	networks, err := calculator.Deaggregate(startStr, endStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	fmt.Printf("Deaggregating %s - %s\n", startStr, endStr)
	fmt.Println(formatter.FormatDeaggregation(networks, format))
}

// handleSplit handles the split mode
func handleSplit(networkStr string, sizeStrs []string, format formatter.OutputFormat) {
	// Parse the network
	var ipStr, maskStr string
	if strings.Contains(networkStr, "/") {
		parts := strings.SplitN(networkStr, "/", 2)
		ipStr = parts[0]
		maskStr = parts[1]
	} else if len(sizeStrs) > 0 {
		ipStr = networkStr
		maskStr = sizeStrs[0]
		sizeStrs = sizeStrs[1:]
	} else {
		fmt.Fprintln(os.Stderr, "Error: No netmask specified")
		os.Exit(1)
	}

	// Parse the sizes
	var sizes []int
	for _, sizeStr := range sizeStrs {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid size: %s\n", sizeStr)
			os.Exit(1)
		}
		sizes = append(sizes, size)
	}

	// Check if it's an IPv6 address
	if strings.Contains(ipStr, ":") {
		fmt.Fprintln(os.Stderr, "Error: IPv6 splitting is not supported yet")
		os.Exit(1)
	}

	// Split the network
	networks, err := calculator.SplitNetwork(ipStr, maskStr, sizes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	fmt.Printf("Splitting %s/%s into subnets\n", ipStr, maskStr)
	fmt.Println(formatter.FormatSplitNetwork(networks, format))
}

// handleNormal handles the normal mode
func handleNormal(args []string, format formatter.OutputFormat) {
	// Parse the IP address and netmask
	var ipStr, maskStr string
	if len(args) == 1 {
		// Check if the argument contains a slash
		if strings.Contains(args[0], "/") {
			parts := strings.SplitN(args[0], "/", 2)
			ipStr = parts[0]
			maskStr = parts[1]
		} else {
			ipStr = args[0]
			// Use default netmask based on IP version
			if strings.Contains(ipStr, ":") {
				maskStr = "64" // Default for IPv6
			} else {
				maskStr = "24" // Default for IPv4
			}
		}
	} else if len(args) >= 2 {
		ipStr = args[0]
		maskStr = args[1]
	}

	// Check if it's an IPv6 address
	if strings.Contains(ipStr, ":") {
		// Calculate IPv6 network
		network, err := calculator.CalculateIPv6Network(ipStr, maskStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print the result
		fmt.Println(formatter.FormatIPv6Network(network, format))
	} else {
		// Calculate IPv4 network
		network, err := calculator.CalculateNetwork(ipStr, maskStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Print the result
		fmt.Println(formatter.FormatIPv4Network(network, format))
	}
} 