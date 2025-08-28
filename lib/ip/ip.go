package ip

import (
	"fmt"
	"net"
	"strconv"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func argumentError(name string, expected int) error {
	return fmt.Errorf("%s: expected %d arguments", name, expected)
}

func isValidIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_valid_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		return lang.BoolValue(ip != nil), nil
	}
	return name, fn
}

func isIPv4() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_ipv4"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return lang.BoolValue(false), nil
		}
		return lang.BoolValue(ip.To4() != nil), nil
	}
	return name, fn
}

func isIPv6() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_ipv6"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return lang.BoolValue(false), nil
		}
		return lang.BoolValue(ip.To4() == nil && ip.To16() != nil), nil
	}
	return name, fn
}

func isPrivateIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_private_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		return lang.BoolValue(ip.IsPrivate()), nil
	}
	return name, fn
}

func isLoopbackIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_loopback_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		return lang.BoolValue(ip.IsLoopback()), nil
	}
	return name, fn
}

func isMulticastIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_multicast_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		return lang.BoolValue(ip.IsMulticast()), nil
	}
	return name, fn
}

func isLinkLocalIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_link_local_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		return lang.BoolValue(ip.IsLinkLocalUnicast()), nil
	}
	return name, fn
}

func cidrMatch() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_match"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, argumentError(name, 2)
		}
		ipStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: ip %w", name, err)
		}
		cidrStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cidr %w", name, err)
		}

		_, network, err := net.ParseCIDR(string(cidrStr))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR '%s': %w", name, string(cidrStr), err)
		}
		ipAddr := net.ParseIP(string(ipStr))
		if ipAddr == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(ipStr))
		}

		return lang.BoolValue(network.Contains(ipAddr)), nil
	}
	return name, fn
}

func cidrContains() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_contains"
	_, cidrMatch := cidrMatch()
	fn := func(args []lang.Value) (lang.Value, error) {
		return cidrMatch(args)
	}
	return name, fn
}

func ipInRange() (string, func([]lang.Value) (lang.Value, error)) {
	name := "ip_in_range"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, argumentError(name, 3)
		}
		ipStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: ip %w", name, err)
		}
		startStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: start ip %w", name, err)
		}
		endStr, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: end ip %w", name, err)
		}

		ipAddr := net.ParseIP(string(ipStr))
		startAddr := net.ParseIP(string(startStr))
		endAddr := net.ParseIP(string(endStr))

		if ipAddr == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(ipStr))
		}
		if startAddr == nil {
			return nil, fmt.Errorf("%s: invalid start IP address '%s'", name, string(startStr))
		}
		if endAddr == nil {
			return nil, fmt.Errorf("%s: invalid end IP address '%s'", name, string(endStr))
		}

		if ipAddr.To4() != nil && startAddr.To4() != nil && endAddr.To4() != nil {
			ipInt := ipv4ToInt(ipAddr.To4())
			startInt := ipv4ToInt(startAddr.To4())
			endInt := ipv4ToInt(endAddr.To4())
			return lang.BoolValue(ipInt >= startInt && ipInt <= endInt), nil
		} else if ipAddr.To16() != nil && startAddr.To16() != nil && endAddr.To16() != nil {
			return lang.BoolValue(compareIPv6(ipAddr.To16(), startAddr.To16()) >= 0 &&
				compareIPv6(ipAddr.To16(), endAddr.To16()) <= 0), nil
		}

		return nil, fmt.Errorf("%s: IP addresses must be of the same type (IPv4 or IPv6)", name)
	}
	return name, fn
}

func cidrNetworkAddress() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_network"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		_, network, err := net.ParseCIDR(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR '%s': %w", name, string(str), err)
		}

		return lang.StringValue(network.IP.String()), nil
	}
	return name, fn
}

func cidrBroadcastAddress() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_broadcast"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		_, network, err := net.ParseCIDR(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR '%s': %w", name, string(str), err)
		}
		ip := network.IP.To4()
		mask := net.IP(network.Mask).To4()
		if ip == nil || mask == nil {
			return nil, fmt.Errorf("%s: IPv6 doesn't have broadcast addresses", name)
		}
		broadcast := make(net.IP, 4)
		for i := range ip {
			broadcast[i] = ip[i] | (^mask[i])
		}
		return lang.StringValue(broadcast.String()), nil
	}
	return name, fn
}

func cidrHostCount() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_host_count"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		_, network, err := net.ParseCIDR(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR '%s': %w", name, string(str), err)
		}
		ones, bits := network.Mask.Size()
		if bits == 32 {
			hostBits := bits - ones
			if hostBits <= 30 {
				count := (1 << hostBits) - 2
				if count < 0 {
					count = 0
				}
				return lang.NumberValue(float64(count)), nil
			}
			return lang.NumberValue(float64((int(1) << hostBits) - 2)), nil
		} else if bits == 128 {
			hostBits := bits - ones
			if hostBits > 50 {
				return lang.NumberValue(-1), nil
			}
			return lang.NumberValue(float64(int(1) << hostBits)), nil
		}
		return lang.NumberValue(0), nil
	}
	return name, fn
}

func cidrSubnets() (string, func([]lang.Value) (lang.Value, error)) {
	name := "cidr_subnets"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, argumentError(name, 2)
		}
		cidrStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: cidr %w", name, err)
		}
		newPrefixLen, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: prefix length %w", name, err)
		}

		_, network, err := net.ParseCIDR(string(cidrStr))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid CIDR '%s': %w", name, string(cidrStr), err)
		}

		prefixLen := int(newPrefixLen)
		ones, bits := network.Mask.Size()
		if prefixLen <= ones || prefixLen > bits {
			return nil, fmt.Errorf("%s: new prefix length %d must be between %d and %d", name, prefixLen, ones+1, bits)
		}
		if bits != 32 {
			return nil, fmt.Errorf("%s: only IPv4 subnets are supported", name)
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
	return name, fn
}

func normalizeIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "normalize_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		return lang.StringValue(ip.String()), nil
	}
	return name, fn
}

func expandIPv6() (string, func([]lang.Value) (lang.Value, error)) {
	name := "expand_ipv6"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		if ip.To4() != nil {
			return nil, fmt.Errorf("%s: address is IPv4, not IPv6", name)
		}
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
	return name, fn
}

func compressIPv6() (string, func([]lang.Value) (lang.Value, error)) {
	name := "compress_ipv6"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}
		if ip.To4() != nil {
			return nil, fmt.Errorf("%s: address is IPv4, not IPv6", name)
		}
		return lang.StringValue(ip.String()), nil
	}
	return name, fn
}

func ipToInt() (string, func([]lang.Value) (lang.Value, error)) {
	name := "ip_to_int"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str)).To4()
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid or non-IPv4 address '%s'", name, string(str))
		}
		return lang.NumberValue(float64(ipv4ToInt(ip))), nil
	}
	return name, fn
}

func intToIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "int_to_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		num, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ipInt := uint32(num)
		ip := intToIPv4(ipInt)
		return lang.StringValue(ip.String()), nil
	}
	return name, fn
}

func reverseIP() (string, func([]lang.Value) (lang.Value, error)) {
	name := "reverse_ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		ip := net.ParseIP(string(str))
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid IP address '%s'", name, string(str))
		}

		if ipv4 := ip.To4(); ipv4 != nil {
			return lang.StringValue(strconv.Itoa(int(ipv4[3])) + "." +
				strconv.Itoa(int(ipv4[2])) + "." +
				strconv.Itoa(int(ipv4[1])) + "." +
				strconv.Itoa(int(ipv4[0])) + ".in-addr.arpa"), nil
		} else {
			ipv6 := ip.To16()
			result := ""
			for i := 15; i >= 0; i-- {
				result += strconv.FormatInt(int64(ipv6[i]&0xf), 16) + "."
				result += strconv.FormatInt(int64(ipv6[i]>>4), 16) + "."
			}
			return lang.StringValue(result + "ip6.arpa"), nil
		}
	}
	return name, fn
}

func isRFC1918() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_rfc1918"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, argumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		ip := net.ParseIP(string(str)).To4()
		if ip == nil {
			return nil, fmt.Errorf("%s: invalid or non-IPv4 address '%s'", name, string(str))
		}

		if ip[0] == 10 {
			return lang.BoolValue(true), nil
		}
		if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
			return lang.BoolValue(true), nil
		}
		if ip[0] == 192 && ip[1] == 168 {
			return lang.BoolValue(true), nil
		}

		return lang.BoolValue(false), nil
	}
	return name, fn
}

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

var IpFunctions = []func() (string, func([]lang.Value) (lang.Value, error)){
	isValidIP,
	isIPv4,
	isIPv6,
	isPrivateIP,
	isLoopbackIP,
	isMulticastIP,
	isLinkLocalIP,
	cidrMatch,
	cidrContains,
	ipInRange,
	cidrNetworkAddress,
	cidrBroadcastAddress,
	cidrHostCount,
	cidrSubnets,
	normalizeIP,
	expandIPv6,
	compressIPv6,
	ipToInt,
	intToIP,
	reverseIP,
	isRFC1918,
}
