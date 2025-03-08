package calculator

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPv4Network represents an IPv4 network with its address and mask
type IPv4Network struct {
	Address    uint32
	Netmask    uint32
	BitCount   int
	Broadcast  uint32
	NetworkID  uint32
	HostMin    uint32
	HostMax    uint32
	HostsCount uint32
	Class      string
}

// ParseIPv4 parses an IPv4 address string into a uint32
func ParseIPv4(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}
	
	// Ensure it's an IPv4 address
	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}
	
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]), nil
}

// ParseNetmask parses a netmask string into a uint32 and bit count
// It accepts CIDR notation (e.g., "24" or "/24") or dotted decimal (e.g., "255.255.255.0")
func ParseNetmask(maskStr string) (uint32, int, error) {
	// Remove leading slash if present
	maskStr = strings.TrimPrefix(maskStr, "/")
	
	// Try to parse as CIDR bit count
	bitCount, err := strconv.Atoi(maskStr)
	if err == nil {
		if bitCount < 0 || bitCount > 32 {
			return 0, 0, fmt.Errorf("invalid bit count: %d (must be between 0 and 32)", bitCount)
		}
		
		// Calculate the netmask from the bit count
		var mask uint32
		for i := 0; i < bitCount; i++ {
			mask |= 1 << (31 - i)
		}
		
		return mask, bitCount, nil
	}
	
	// Try to parse as dotted decimal
	mask, err := ParseIPv4(maskStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid netmask: %s", maskStr)
	}
	
	// Validate the netmask (must be contiguous 1s followed by contiguous 0s)
	if !isValidNetmask(mask) {
		return 0, 0, fmt.Errorf("invalid netmask: %s (not contiguous)", maskStr)
	}
	
	// Count the bits
	bitCount = 0
	for i := 0; i < 32; i++ {
		if (mask & (1 << (31 - i))) != 0 {
			bitCount++
		} else {
			break
		}
	}
	
	return mask, bitCount, nil
}

// isValidNetmask checks if a netmask is valid (contiguous 1s followed by contiguous 0s)
func isValidNetmask(mask uint32) bool {
	// Find the first 0 bit
	var foundZero bool
	for i := 31; i >= 0; i-- {
		bit := (mask & (1 << i)) != 0
		if !bit {
			foundZero = true
		} else if foundZero {
			// Found a 1 after a 0, invalid mask
			return false
		}
	}
	return true
}

// IPToString converts a uint32 IP address to a dotted decimal string
func IPToString(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip>>24)&0xFF,
		(ip>>16)&0xFF,
		(ip>>8)&0xFF,
		ip&0xFF)
}

// GetClass returns the class of an IPv4 address
func GetClass(ip uint32) string {
	// Class A: 0.0.0.0 to 127.255.255.255
	if (ip & 0x80000000) == 0 {
		return "A"
	}
	// Class B: 128.0.0.0 to 191.255.255.255
	if (ip & 0xC0000000) == 0x80000000 {
		return "B"
	}
	// Class C: 192.0.0.0 to 223.255.255.255
	if (ip & 0xE0000000) == 0xC0000000 {
		return "C"
	}
	// Class D: 224.0.0.0 to 239.255.255.255
	if (ip & 0xF0000000) == 0xE0000000 {
		return "D"
	}
	// Class E: 240.0.0.0 to 255.255.255.255
	return "E"
}

// GetClassBits returns the natural bit count for the class
func GetClassBits(class string) int {
	switch class {
	case "A":
		return 8
	case "B":
		return 16
	case "C":
		return 24
	case "D", "E":
		return 4
	default:
		return 0
	}
}

// CalculateNetwork calculates network details from an IP address and netmask
func CalculateNetwork(ipStr, maskStr string) (*IPv4Network, error) {
	ip, err := ParseIPv4(ipStr)
	if err != nil {
		return nil, err
	}
	
	mask, bitCount, err := ParseNetmask(maskStr)
	if err != nil {
		return nil, err
	}
	
	network := &IPv4Network{
		Address:   ip,
		Netmask:   mask,
		BitCount:  bitCount,
		NetworkID: ip & mask,
	}
	
	network.Broadcast = network.NetworkID | ^mask
	
	// Special case for /31 and /32 networks
	if bitCount == 31 {
		network.HostMin = network.NetworkID
		network.HostMax = network.Broadcast
		network.HostsCount = 2
	} else if bitCount == 32 {
		network.HostMin = network.NetworkID
		network.HostMax = network.NetworkID
		network.HostsCount = 1
	} else {
		network.HostMin = network.NetworkID + 1
		network.HostMax = network.Broadcast - 1
		network.HostsCount = (1 << (32 - bitCount)) - 2
	}
	
	network.Class = GetClass(ip)
	
	return network, nil
}

// GetWildcardMask returns the wildcard mask (inverse of netmask)
func GetWildcardMask(netmask uint32) uint32 {
	return ^netmask
}

// IsPrivate checks if an IP address is in a private range
func IsPrivate(ip uint32) bool {
	// 10.0.0.0/8
	if (ip & 0xFF000000) == 0x0A000000 {
		return true
	}
	// 172.16.0.0/12
	if (ip & 0xFFF00000) == 0xAC100000 {
		return true
	}
	// 192.168.0.0/16
	if (ip & 0xFFFF0000) == 0xC0A80000 {
		return true
	}
	return false
}

// Deaggregate returns a list of CIDR blocks that cover the range from start to end
func Deaggregate(startStr, endStr string) ([]string, error) {
	start, err := ParseIPv4(startStr)
	if err != nil {
		return nil, err
	}
	
	end, err := ParseIPv4(endStr)
	if err != nil {
		return nil, err
	}
	
	if start > end {
		return nil, errors.New("start address must be less than or equal to end address")
	}
	
	var result []string
	
	for start <= end {
		// Find the largest block that fits
		var prefix int
		for prefix = 0; prefix < 32; prefix++ {
			mask := uint32(0xFFFFFFFF) << (32 - prefix)
			networkID := start & mask
			
			// Check if this block fits within our range
			blockEnd := networkID | ^mask
			if blockEnd > end {
				continue
			}
			
			// Check if this is a valid block starting at 'start'
			if networkID == start {
				break
			}
		}
		
		// Add the block to our result
		result = append(result, fmt.Sprintf("%s/%d", IPToString(start), prefix))
		
		// Move to the next block
		start += 1 << (32 - prefix)
	}
	
	return result, nil
}

// SplitNetwork splits a network into subnets of specified sizes
func SplitNetwork(networkStr, maskStr string, sizes []int) ([]string, error) {
	network, err := CalculateNetwork(networkStr, maskStr)
	if err != nil {
		return nil, err
	}
	
	// Calculate total hosts in the network
	totalHosts := uint32(1) << (32 - uint32(network.BitCount))
	
	// Convert sizes to host counts and check if they fit
	var hostCounts []uint32
	var totalRequired uint32
	
	for _, size := range sizes {
		// Calculate required hosts for this subnet (including network and broadcast)
		var hostsNeeded uint32 = 4 // Start with at least 4 hosts (network, broadcast, and 2 usable)
		for hostsNeeded-2 < uint32(size) {
			hostsNeeded *= 2
		}
		
		hostCounts = append(hostCounts, hostsNeeded)
		totalRequired += hostsNeeded
	}
	
	if totalRequired > totalHosts {
		return nil, fmt.Errorf("requested subnet sizes exceed available space (%d > %d)", totalRequired, totalHosts)
	}
	
	// Allocate subnets
	var result []string
	currentIP := network.NetworkID
	
	for _, hostCount := range hostCounts {
		// Calculate prefix for this subnet
		prefix := 32
		for hostCount > 1 {
			hostCount >>= 1
			prefix--
		}
		
		// Create a subnet object to get proper formatting
		subnetNetwork, err := CalculateNetwork(IPToString(currentIP), fmt.Sprintf("%d", prefix))
		if err != nil {
			return nil, err
		}
		
		// Add subnet to result with CIDR notation
		result = append(result, fmt.Sprintf("%s/%d", IPToString(subnetNetwork.NetworkID), prefix))
		
		// Move to next subnet
		currentIP += 1 << (32 - prefix)
	}
	
	return result, nil
}

// FormatBinary returns the binary representation of an IP address
func FormatBinary(ip uint32) string {
	var parts []string
	for i := 0; i < 4; i++ {
		octet := (ip >> (24 - i*8)) & 0xFF
		binary := fmt.Sprintf("%08b", octet)
		parts = append(parts, binary)
	}
	return strings.Join(parts, ".")
} 