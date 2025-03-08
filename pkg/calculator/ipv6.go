package calculator

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// IPv6Network represents an IPv6 network with its address and prefix
type IPv6Network struct {
	Address     *big.Int
	PrefixLen   int
	NetworkID   *big.Int
	NetworkMask *big.Int
}

// ParseIPv6 parses an IPv6 address string into a big.Int
func ParseIPv6(ipStr string) (*big.Int, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}
	
	// Ensure it's an IPv6 address
	ip = ip.To16()
	if ip == nil {
		return nil, fmt.Errorf("not an IPv6 address: %s", ipStr)
	}
	
	// Convert to big.Int
	ipInt := new(big.Int)
	ipInt.SetBytes(ip)
	
	return ipInt, nil
}

// ParseIPv6Prefix parses an IPv6 prefix length
func ParseIPv6Prefix(prefixStr string) (int, error) {
	// Remove leading slash if present
	prefixStr = strings.TrimPrefix(prefixStr, "/")
	
	// Parse as integer
	prefix, err := strconv.Atoi(prefixStr)
	if err != nil {
		return 0, fmt.Errorf("invalid prefix length: %s", prefixStr)
	}
	
	// Validate range
	if prefix < 0 || prefix > 128 {
		return 0, fmt.Errorf("invalid prefix length: %d (must be between 0 and 128)", prefix)
	}
	
	return prefix, nil
}

// IPv6ToNetworkID calculates the network ID for an IPv6 address with a given prefix length
func IPv6ToNetworkID(ip *big.Int, prefixLen int) (*big.Int, error) {
	// Create a mask with prefixLen 1's followed by (128-prefixLen) 0's
	mask := new(big.Int)
	mask.Lsh(big.NewInt(1), 128)    // mask = 2^128
	mask.Sub(mask, big.NewInt(1))   // mask = 2^128 - 1 (all 1's)
	
	// Shift left to remove the host bits
	if prefixLen < 0 {
		return nil, fmt.Errorf("prefix length cannot be negative: %d", prefixLen)
	}
	if prefixLen > 128 {
		return nil, fmt.Errorf("prefix length cannot exceed 128: %d", prefixLen)
	}
	
	if prefixLen < 128 {
		shiftBits := 128 - prefixLen // Calculate shift bits as int first
		// Ensure shiftBits is positive before converting to uint
		if shiftBits < 0 {
			return nil, fmt.Errorf("invalid shift bits calculation: %d", shiftBits)
		}
		// Now it's safe to convert to uint
		uShiftBits := uint(shiftBits)
		mask.Rsh(mask, uShiftBits)
		mask.Lsh(mask, uShiftBits)
	}
	
	// Apply the mask to the IP
	result := new(big.Int)
	result.And(ip, mask)
	
	return result, nil
}

// IPv6ToString converts a big.Int to an IPv6 address string
func IPv6ToString(ipInt *big.Int) string {
	// Convert to 16-byte array
	ipBytes := make([]byte, 16)
	bytes := ipInt.Bytes()
	
	// Pad with leading zeros if needed
	copy(ipBytes[16-len(bytes):], bytes)
	
	// Convert to net.IP and format
	ip := net.IP(ipBytes)
	return ip.String()
}

// CalculateIPv6Network calculates network details from an IPv6 address and prefix
func CalculateIPv6Network(ipStr, prefixStr string) (*IPv6Network, error) {
	ip, err := ParseIPv6(ipStr)
	if err != nil {
		return nil, err
	}
	
	prefix, err := ParseIPv6Prefix(prefixStr)
	if err != nil {
		return nil, err
	}
	
	// Calculate network mask
	networkMask := new(big.Int)
	networkMask.Lsh(big.NewInt(1), 128)    // mask = 2^128
	networkMask.Sub(networkMask, big.NewInt(1))   // mask = 2^128 - 1 (all 1's)
	
	// Shift left to remove the host bits
	if prefix < 0 {
		return nil, fmt.Errorf("prefix length cannot be negative: %d", prefix)
	}
	if prefix > 128 {
		return nil, fmt.Errorf("prefix length cannot exceed 128: %d", prefix)
	}
	
	if prefix < 128 {
		shiftBits := 128 - prefix // Calculate shift bits as int first
		// Ensure shiftBits is positive before converting to uint
		if shiftBits < 0 {
			return nil, fmt.Errorf("invalid shift bits calculation: %d", shiftBits)
		}
		// Now it's safe to convert to uint
		uShiftBits := uint(shiftBits)
		networkMask.Rsh(networkMask, uShiftBits)
		networkMask.Lsh(networkMask, uShiftBits)
	}
	
	// Calculate network ID
	networkID := new(big.Int)
	networkID.And(ip, networkMask)
	
	return &IPv6Network{
		Address:     ip,
		PrefixLen:   prefix,
		NetworkID:   networkID,
		NetworkMask: networkMask,
	}, nil
}

// FormatIPv6Binary returns the binary representation of an IPv6 address
func FormatIPv6Binary(ip *big.Int) string {
	// Convert to 16-byte array
	ipBytes := make([]byte, 16)
	bytes := ip.Bytes()
	
	// Pad with leading zeros if needed
	copy(ipBytes[16-len(bytes):], bytes)
	
	var parts []string
	for _, b := range ipBytes {
		parts = append(parts, fmt.Sprintf("%08b", b))
	}
	
	// Group into 16-bit chunks for readability
	var result []string
	for i := 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			result = append(result, parts[i]+parts[i+1])
		} else {
			result = append(result, parts[i])
		}
	}
	
	return strings.Join(result, ":")
} 