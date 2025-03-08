package formatter

import (
	"fmt"
	"strings"

	"github.com/neontowel/ipcalc-go/pkg/calculator"
)

// OutputFormat defines the format for the output
type OutputFormat struct {
	UseColor  bool
	UseHTML   bool
	UseBinary bool
}

// ColorCodes for terminal output
type ColorCodes struct {
	Reset    string
	Address  string
	Netmask  string
	Binary   string
	Class    string
	Subnet   string
	Error    string
	Wildcard string
}

// DefaultColors returns the default color codes
func DefaultColors() ColorCodes {
	return ColorCodes{
		Reset:    "\033[0m",
		Address:  "\033[34m", // Blue
		Netmask:  "\033[31m", // Red
		Binary:   "\033[33m", // Yellow
		Class:    "\033[35m", // Magenta
		Subnet:   "\033[32m", // Green
		Error:    "\033[31m", // Red
		Wildcard: "\033[36m", // Cyan
	}
}

// HTMLColors returns HTML color codes
func HTMLColors() ColorCodes {
	return ColorCodes{
		Reset:    "</font>",
		Address:  "<font color=\"#0000ff\">",
		Netmask:  "<font color=\"#ff0000\">",
		Binary:   "<font color=\"#909090\">",
		Class:    "<font color=\"#009900\">",
		Subnet:   "<font color=\"#663366\">",
		Error:    "<font color=\"#ff0000\">",
		Wildcard: "<font color=\"#00cccc\">",
	}
}

// NoColors returns empty color codes
func NoColors() ColorCodes {
	return ColorCodes{}
}

// FormatIPv4Network formats an IPv4Network for display
func FormatIPv4Network(network *calculator.IPv4Network, format OutputFormat) string {
	var colors ColorCodes
	var lineBreak string

	if format.UseHTML {
		colors = HTMLColors()
		lineBreak = "<br>\n"
	} else if format.UseColor {
		colors = DefaultColors()
		lineBreak = "\n"
	} else {
		colors = NoColors()
		lineBreak = "\n"
	}

	var result strings.Builder

	// Address line
	result.WriteString(fmt.Sprintf("Address:   %s%s%s", 
		colors.Address, 
		calculator.IPToString(network.Address), 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("          %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(network.Address), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// Netmask line
	result.WriteString(fmt.Sprintf("Netmask:   %s%s = %d%s", 
		colors.Netmask, 
		calculator.IPToString(network.Netmask), 
		network.BitCount, 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("   %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(network.Netmask), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// Wildcard line
	wildcard := calculator.GetWildcardMask(network.Netmask)
	result.WriteString(fmt.Sprintf("Wildcard:  %s%s%s", 
		colors.Wildcard, 
		calculator.IPToString(wildcard), 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("            %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(wildcard), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	result.WriteString("=>" + lineBreak)

	// Network line
	result.WriteString(fmt.Sprintf("Network:   %s%s/%d%s", 
		colors.Subnet, 
		calculator.IPToString(network.NetworkID), 
		network.BitCount, 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("       %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(network.NetworkID), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// HostMin line
	result.WriteString(fmt.Sprintf("HostMin:   %s%s%s", 
		colors.Subnet, 
		calculator.IPToString(network.HostMin), 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("          %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(network.HostMin), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// HostMax line
	result.WriteString(fmt.Sprintf("HostMax:   %s%s%s", 
		colors.Subnet, 
		calculator.IPToString(network.HostMax), 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("          %s%s%s", 
			colors.Binary, 
			calculator.FormatBinary(network.HostMax), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// Broadcast line (only for masks < 31)
	if network.BitCount < 31 {
		result.WriteString(fmt.Sprintf("Broadcast: %s%s%s", 
			colors.Subnet, 
			calculator.IPToString(network.Broadcast), 
			colors.Reset))
		
		if format.UseBinary {
			result.WriteString(fmt.Sprintf("          %s%s%s", 
				colors.Binary, 
				calculator.FormatBinary(network.Broadcast), 
				colors.Reset))
		}
		result.WriteString(lineBreak)
	}

	// Hosts/Net line
	result.WriteString(fmt.Sprintf("Hosts/Net: %s%d%s", 
		colors.Subnet, 
		network.HostsCount, 
		colors.Reset))

	// Class info
	classInfo := fmt.Sprintf("Class %s", network.Class)
	if calculator.IsPrivate(network.Address) {
		classInfo += ", Private Internet"
	}
	result.WriteString(fmt.Sprintf("                   %s%s%s", 
		colors.Class, 
		classInfo, 
		colors.Reset))
	
	return result.String()
}

// FormatIPv6Network formats an IPv6Network for display
func FormatIPv6Network(network *calculator.IPv6Network, format OutputFormat) string {
	var colors ColorCodes
	var lineBreak string

	if format.UseHTML {
		colors = HTMLColors()
		lineBreak = "<br>\n"
	} else if format.UseColor {
		colors = DefaultColors()
		lineBreak = "\n"
	} else {
		colors = NoColors()
		lineBreak = "\n"
	}

	var result strings.Builder

	// Address line
	result.WriteString(fmt.Sprintf("Address: %s%s%s", 
		colors.Address, 
		calculator.IPv6ToString(network.Address), 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("     %s%s%s", 
			colors.Binary, 
			calculator.FormatIPv6Binary(network.Address), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// Netmask line
	result.WriteString(fmt.Sprintf("Netmask: %s%d%s", 
		colors.Netmask, 
		network.PrefixLen, 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("                                      %s%s%s", 
			colors.Binary, 
			calculator.FormatIPv6Binary(network.NetworkMask), 
			colors.Reset))
	}
	result.WriteString(lineBreak)

	// Prefix line
	result.WriteString(fmt.Sprintf("Prefix:  %s%s/%d%s", 
		colors.Subnet, 
		calculator.IPv6ToString(network.NetworkID), 
		network.PrefixLen, 
		colors.Reset))
	
	if format.UseBinary {
		result.WriteString(fmt.Sprintf("                     %s%s%s", 
			colors.Binary, 
			calculator.FormatIPv6Binary(network.NetworkID), 
			colors.Reset))
	}
	
	return result.String()
}

// FormatDeaggregation formats the results of a deaggregation
func FormatDeaggregation(networks []string, format OutputFormat) string {
	var colors ColorCodes
	var lineBreak string

	if format.UseHTML {
		colors = HTMLColors()
		lineBreak = "<br>\n"
	} else if format.UseColor {
		colors = DefaultColors()
		lineBreak = "\n"
	} else {
		colors = NoColors()
		lineBreak = "\n"
	}

	var result strings.Builder
	
	for _, network := range networks {
		result.WriteString(fmt.Sprintf("%s%s%s%s", 
			colors.Subnet, 
			network, 
			colors.Reset,
			lineBreak))
	}
	
	return result.String()
}

// FormatSplitNetwork formats the results of a network split
func FormatSplitNetwork(networks []string, format OutputFormat) string {
	var colors ColorCodes
	var lineBreak string

	if format.UseHTML {
		colors = HTMLColors()
		lineBreak = "<br>\n"
	} else if format.UseColor {
		colors = DefaultColors()
		lineBreak = "\n"
	} else {
		colors = NoColors()
		lineBreak = "\n"
	}

	var result strings.Builder
	
	for i, network := range networks {
		result.WriteString(fmt.Sprintf("Subnet %d: %s%s%s%s", 
			i+1,
			colors.Subnet, 
			network, 
			colors.Reset,
			lineBreak))
	}
	
	return result.String()
}

// FormatHTMLHeader returns the HTML header
func FormatHTMLHeader() string {
	return `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
<head>
<meta HTTP-EQUIV="content-type" CONTENT="text/html; charset=UTF-8">
<title>IP Calculator</title>
</head>
<body>
`
}

// FormatHTMLFooter returns the HTML footer
func FormatHTMLFooter() string {
	return `</body>
</html>
`
} 