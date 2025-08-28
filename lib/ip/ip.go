package ip

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Network/IP Functions
// These functions help with IP address validation, CIDR operations, and network analysis

func isValidIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_valid_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_valid_ip: %w", err)
	}
	ip := net.ParseIP(string(str))
	return lang.BoolValue(ip != nil), nil
}

func isIPv4(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_ipv4: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_ipv4: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return lang.BoolValue(false), nil
	}
	return lang.BoolValue(ip.To4() != nil), nil
}

func isIPv6(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_ipv6: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_ipv6: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return lang.BoolValue(false), nil
	}
	return lang.BoolValue(ip.To4() == nil && ip.To16() != nil), nil
}

func isPrivateIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_private_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_private_ip: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("is_private_ip: invalid IP address '%s'", string(str))
	}
	return lang.BoolValue(ip.IsPrivate()), nil
}

func isLoopbackIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_loopback_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_loopback_ip: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("is_loopback_ip: invalid IP address '%s'", string(str))
	}
	return lang.BoolValue(ip.IsLoopback()), nil
}

func isMulticastIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_multicast_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_multicast_ip: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("is_multicast_ip: invalid IP address '%s'", string(str))
	}
	return lang.BoolValue(ip.IsMulticast()), nil
}

func isLinkLocalIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_link_local_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_link_local_ip: %w", err)
	}
	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("is_link_local_ip: invalid IP address '%s'", string(str))
	}
	return lang.BoolValue(ip.IsLinkLocalUnicast()), nil
}

func cidrMatch(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("cidr_match: expected 2 arguments (ip, cidr)")
	}
	ipStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("cidr_match: ip %w", err)
	}
	cidrStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("cidr_match: cidr %w", err)
	}

	_, network, err := net.ParseCIDR(string(cidrStr))
	if err != nil {
		return nil, fmt.Errorf("cidr_match: invalid CIDR '%s': %w", string(cidrStr), err)
	}

	ipAddr := net.ParseIP(string(ipStr))
	if ipAddr == nil {
		return nil, fmt.Errorf("cidr_match: invalid IP address '%s'", string(ipStr))
	}

	return lang.BoolValue(network.Contains(ipAddr)), nil
}

func cidrContains(args []lang.Value) (lang.Value, error) {
	// Alias for cidr_match for clarity
	return cidrMatch(args)
}

func ipInRange(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("ip_in_range: expected 3 arguments (ip, start_ip, end_ip)")
	}
	ipStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("ip_in_range: ip %w", err)
	}
	startStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("ip_in_range: start ip %w", err)
	}
	endStr, err := lib.ToString(args[2])
	if err != nil {
		return nil, fmt.Errorf("ip_in_range: end ip %w", err)
	}

	ipAddr := net.ParseIP(string(ipStr))
	startAddr := net.ParseIP(string(startStr))
	endAddr := net.ParseIP(string(endStr))

	if ipAddr == nil {
		return nil, fmt.Errorf("ip_in_range: invalid IP address '%s'", string(ipStr))
	}
	if startAddr == nil {
		return nil, fmt.Errorf("ip_in_range: invalid start IP address '%s'", string(startStr))
	}
	if endAddr == nil {
		return nil, fmt.Errorf("ip_in_range: invalid end IP address '%s'", string(endStr))
	}

	// Convert to comparable format
	if ipAddr.To4() != nil && startAddr.To4() != nil && endAddr.To4() != nil {
		// IPv4 comparison
		ipInt := ipv4ToInt(ipAddr.To4())
		startInt := ipv4ToInt(startAddr.To4())
		endInt := ipv4ToInt(endAddr.To4())
		return lang.BoolValue(ipInt >= startInt && ipInt <= endInt), nil
	} else if ipAddr.To16() != nil && startAddr.To16() != nil && endAddr.To16() != nil {
		// IPv6 comparison - use byte comparison
		return lang.BoolValue(compareIPv6(ipAddr.To16(), startAddr.To16()) >= 0 &&
			compareIPv6(ipAddr.To16(), endAddr.To16()) <= 0), nil
	}

	return nil, errors.New("ip_in_range: IP addresses must be of the same type (IPv4 or IPv6)")
}

func cidrNetworkAddress(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cidr_network: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("cidr_network: %w", err)
	}

	_, network, err := net.ParseCIDR(string(str))
	if err != nil {
		return nil, fmt.Errorf("cidr_network: invalid CIDR '%s': %w", string(str), err)
	}

	return lang.StringValue(network.IP.String()), nil
}

func cidrBroadcastAddress(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cidr_broadcast: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("cidr_broadcast: %w", err)
	}

	_, network, err := net.ParseCIDR(string(str))
	if err != nil {
		return nil, fmt.Errorf("cidr_broadcast: invalid CIDR '%s': %w", string(str), err)
	}

	// Calculate broadcast address
	ip := network.IP.To4()
	mask := net.IP(network.Mask).To4()

	if ip == nil || mask == nil {
		return nil, errors.New("cidr_broadcast: IPv6 doesn't have broadcast addresses")
	}

	broadcast := make(net.IP, 4)
	for i := range ip {
		broadcast[i] = ip[i] | (^mask[i])
	}

	return lang.StringValue(broadcast.String()), nil
}

func cidrHostCount(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cidr_host_count: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("cidr_host_count: %w", err)
	}

	_, network, err := net.ParseCIDR(string(str))
	if err != nil {
		return nil, fmt.Errorf("cidr_host_count: invalid CIDR '%s': %w", string(str), err)
	}

	ones, bits := network.Mask.Size()
	if bits == 32 {
		// IPv4
		hostBits := bits - ones
		if hostBits <= 30 {
			count := (1 << hostBits) - 2 // Subtract network and broadcast
			if count < 0 {
				count = 0
			}
			return lang.NumberValue(float64(count)), nil
		}
		return lang.NumberValue(float64((int(1) << hostBits) - 2)), nil
	} else if bits == 128 {
		// IPv6 - return approximation for large networks
		hostBits := bits - ones
		if hostBits > 50 {
			return lang.NumberValue(-1), nil // Too large to represent
		}
		return lang.NumberValue(float64(int(1) << hostBits)), nil
	}

	return lang.NumberValue(0), nil
}

func cidrSubnets(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("cidr_subnets: expected 2 arguments (cidr, new_prefix_length)")
	}
	cidrStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("cidr_subnets: cidr %w", err)
	}
	newPrefixLen, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("cidr_subnets: prefix length %w", err)
	}

	_, network, err := net.ParseCIDR(string(cidrStr))
	if err != nil {
		return nil, fmt.Errorf("cidr_subnets: invalid CIDR '%s': %w", string(cidrStr), err)
	}

	prefixLen := int(newPrefixLen)
	ones, bits := network.Mask.Size()
	if prefixLen <= ones || prefixLen > bits {
		return nil, fmt.Errorf("cidr_subnets: new prefix length %d must be between %d and %d", prefixLen, ones+1, bits)
	}

	// Only handle IPv4 for simplicity
	if bits != 32 {
		return nil, errors.New("cidr_subnets: only IPv4 subnets are supported")
	}

	subnets := make(lang.ListValue, 0)
	subnetCount := 1 << (prefixLen - ones)
	hostSize := 1 << (32 - prefixLen)

	baseIP := ipv4ToInt(network.IP.To4())

	for i := 0; i < subnetCount; i++ {
		subnetIP := intToIPv4(baseIP + uint32(i*hostSize))
		subnet := subnetIP.String() + "/" + strconv.Itoa(prefixLen)
		subnets = append(subnets, lang.StringValue(subnet))
	}

	return subnets, nil
}

func normalizeIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("normalize_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("normalize_ip: %w", err)
	}

	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("normalize_ip: invalid IP address '%s'", string(str))
	}

	return lang.StringValue(ip.String()), nil
}

func expandIPv6(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("expand_ipv6: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("expand_ipv6: %w", err)
	}

	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("expand_ipv6: invalid IP address '%s'", string(str))
	}
	if ip.To4() != nil {
		return nil, errors.New("expand_ipv6: address is IPv4, not IPv6")
	}

	// Format as full IPv6 address
	ipv6 := ip.To16()
	result := ""
	for i := 0; i < 16; i += 2 {
		if i > 0 {
			result += ":"
		}
		result += fmt.Sprintf("%04x", int(ipv6[i])<<8+int(ipv6[i+1]))
	}

	return lang.StringValue(result), nil
}

func compressIPv6(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("compress_ipv6: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("compress_ipv6: %w", err)
	}

	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("compress_ipv6: invalid IP address '%s'", string(str))
	}
	if ip.To4() != nil {
		return nil, errors.New("compress_ipv6: address is IPv4, not IPv6")
	}

	return lang.StringValue(ip.String()), nil
}

func ipToInt(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("ip_to_int: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("ip_to_int: %w", err)
	}

	ip := net.ParseIP(string(str)).To4()
	if ip == nil {
		return nil, fmt.Errorf("ip_to_int: invalid or non-IPv4 address '%s'", string(str))
	}

	return lang.NumberValue(float64(ipv4ToInt(ip))), nil
}

func intToIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("int_to_ip: expected 1 argument")
	}
	num, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("int_to_ip: %w", err)
	}

	ipInt := uint32(num)
	ip := intToIPv4(ipInt)
	return lang.StringValue(ip.String()), nil
}

func reverseIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("reverse_ip: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("reverse_ip: %w", err)
	}

	ip := net.ParseIP(string(str))
	if ip == nil {
		return nil, fmt.Errorf("reverse_ip: invalid IP address '%s'", string(str))
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		// IPv4 reverse DNS format
		return lang.StringValue(strconv.Itoa(int(ipv4[3])) + "." +
			strconv.Itoa(int(ipv4[2])) + "." +
			strconv.Itoa(int(ipv4[1])) + "." +
			strconv.Itoa(int(ipv4[0])) + ".in-addr.arpa"), nil
	} else {
		// IPv6 reverse DNS format
		ipv6 := ip.To16()
		result := ""
		for i := 15; i >= 0; i-- {
			result += strconv.FormatInt(int64(ipv6[i]&0xf), 16) + "."
			result += strconv.FormatInt(int64(ipv6[i]>>4), 16) + "."
		}
		return lang.StringValue(result + "ip6.arpa"), nil
	}
}

func isRFC1918(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_rfc1918: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_rfc1918: %w", err)
	}

	ip := net.ParseIP(string(str)).To4()
	if ip == nil {
		return nil, fmt.Errorf("is_rfc1918: invalid or non-IPv4 address '%s'", string(str))
	}

	// Check RFC 1918 private address ranges
	// 10.0.0.0/8
	if ip[0] == 10 {
		return lang.BoolValue(true), nil
	}
	// 172.16.0.0/12
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return lang.BoolValue(true), nil
	}
	// 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return lang.BoolValue(true), nil
	}

	return lang.BoolValue(false), nil
}

// Helper functions
func ipv4ToInt(ip net.IP) uint32 {
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

func intToIPv4(ipInt uint32) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}

func compareIPv6(a, b net.IP) int {
	for i := 0; i < 16; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// Functions that would be in the BuiltinFunctions map:
var NetworkFunctions = map[string]func([]lang.Value) (lang.Value, error){
	"is_valid_ip":      isValidIP,
	"is_ipv4":          isIPv4,
	"is_ipv6":          isIPv6,
	"is_private_ip":    isPrivateIP,
	"is_loopback_ip":   isLoopbackIP,
	"is_multicast_ip":  isMulticastIP,
	"is_link_local_ip": isLinkLocalIP,
	"cidr_match":       cidrMatch,
	"cidr_contains":    cidrContains,
	"ip_in_range":      ipInRange,
	"cidr_network":     cidrNetworkAddress,
	"cidr_broadcast":   cidrBroadcastAddress,
	"cidr_host_count":  cidrHostCount,
	"cidr_subnets":     cidrSubnets,
	"normalize_ip":     normalizeIP,
	"expand_ipv6":      expandIPv6,
	"compress_ipv6":    compressIPv6,
	"ip_to_int":        ipToInt,
	"int_to_ip":        intToIP,
	"reverse_ip":       reverseIP,
	"is_rfc1918":       isRFC1918,
}
