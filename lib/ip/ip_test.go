package ip

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestIsValidIP(t *testing.T) {
	_, fn := isValidIP()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"valid IPv4", "192.168.1.1", true, false},
		{"valid IPv6", "2001:db8::1", true, false},
		{"invalid IP", "256.256.256.256", false, false},
		{"invalid string", "not-an-ip", false, false},
		{"empty string", "", false, false},
		{"localhost", "127.0.0.1", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	_, fn := isIPv4()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid IPv4", "192.168.1.1", true},
		{"valid IPv6", "2001:db8::1", false},
		{"invalid IP", "invalid", false},
		{"localhost IPv4", "127.0.0.1", true},
		{"localhost IPv6", "::1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	_, fn := isIPv6()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid IPv6", "2001:db8::1", true},
		{"valid IPv4", "192.168.1.1", false},
		{"localhost IPv6", "::1", true},
		{"invalid IP", "invalid", false},
		{"IPv6 full", "2001:0db8:0000:0000:0000:8a2e:0370:7334", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	_, fn := isPrivateIP()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"private IPv4 10.x", "10.0.0.1", true, false},
		{"private IPv4 192.168.x", "192.168.1.1", true, false},
		{"private IPv4 172.x", "172.16.0.1", true, false},
		{"public IPv4", "8.8.8.8", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsLoopbackIP(t *testing.T) {
	_, fn := isLoopbackIP()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"IPv4 loopback", "127.0.0.1", true, false},
		{"IPv6 loopback", "::1", true, false},
		{"public IP", "8.8.8.8", false, false},
		{"invalid IP", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsMulticastIP(t *testing.T) {
	_, fn := isMulticastIP()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"IPv4 multicast", "224.0.0.1", true, false},
		{"IPv6 multicast", "ff02::1", true, false},
		{"unicast IP", "192.168.1.1", false, false},
		{"invalid IP", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsLinkLocalIP(t *testing.T) {
	_, fn := isLinkLocalIP()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"IPv4 link local", "169.254.1.1", true, false},
		{"IPv6 link local", "fe80::1", true, false},
		{"regular IP", "192.168.1.1", false, false},
		{"invalid IP", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestCidrMatch(t *testing.T) {
	_, fn := cidrMatch()

	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
		hasError bool
	}{
		{"IP in subnet", "192.168.1.100", "192.168.1.0/24", true, false},
		{"IP not in subnet", "192.168.2.100", "192.168.1.0/24", false, false},
		{"exact match", "192.168.1.1", "192.168.1.1/32", true, false},
		{"invalid CIDR", "192.168.1.1", "invalid", false, true},
		{"invalid IP", "invalid", "192.168.1.0/24", false, true},
		{"wrong arg count", "", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			if tt.name == "wrong arg count" {
				args = []lang.Value{lang.StringValue(tt.ip)}
			} else {
				args = []lang.Value{lang.StringValue(tt.ip), lang.StringValue(tt.cidr)}
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIPInRange(t *testing.T) {
	_, fn := ipInRange()

	tests := []struct {
		name     string
		ip       string
		start    string
		end      string
		expected bool
		hasError bool
	}{
		{"IP in range", "192.168.1.100", "192.168.1.1", "192.168.1.200", true, false},
		{"IP not in range", "192.168.2.100", "192.168.1.1", "192.168.1.200", false, false},
		{"IP at start", "192.168.1.1", "192.168.1.1", "192.168.1.200", true, false},
		{"IP at end", "192.168.1.200", "192.168.1.1", "192.168.1.200", true, false},
		{"IPv6 range", "2001:db8::100", "2001:db8::1", "2001:db8::200", true, false},
		{"mixed IP types", "192.168.1.100", "2001:db8::1", "2001:db8::200", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.ip), lang.StringValue(tt.start), lang.StringValue(tt.end)}
			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestCidrNetworkAddress(t *testing.T) {
	_, fn := cidrNetworkAddress()

	tests := []struct {
		name     string
		cidr     string
		expected string
		hasError bool
	}{
		{"IPv4 network", "192.168.1.100/24", "192.168.1.0", false},
		{"IPv6 network", "2001:db8::100/64", "2001:db8::", false},
		{"invalid CIDR", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.cidr)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestCidrBroadcastAddress(t *testing.T) {
	_, fn := cidrBroadcastAddress()

	tests := []struct {
		name     string
		cidr     string
		expected string
		hasError bool
	}{
		{"IPv4 broadcast", "192.168.1.0/24", "192.168.1.255", false},
		{"IPv4 /30", "192.168.1.0/30", "192.168.1.3", false},
		{"IPv6 no broadcast", "2001:db8::/64", "", true},
		{"invalid CIDR", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.cidr)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestCidrHostCount(t *testing.T) {
	_, fn := cidrHostCount()

	tests := []struct {
		name     string
		cidr     string
		expected float64
		hasError bool
	}{
		{"IPv4 /24", "192.168.1.0/24", 254, false},
		{"IPv4 /30", "192.168.1.0/30", 2, false},
		{"IPv4 /31", "192.168.1.0/31", 0, false},
		//{"IPv6 /64", "2001:db8::/64", float64(1 << 64), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.cidr)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			actual := float64(result.(lang.NumberValue))
			if tt.expected == float64(1<<64) {
				if actual < 0 {
					t.Errorf("Expected large positive number for IPv6, got %f", actual)
				}
			} else if actual != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, actual)
			}
		})
	}
}

func TestCidrSubnets(t *testing.T) {
	_, fn := cidrSubnets()

	tests := []struct {
		name          string
		cidr          string
		newPrefixLen  float64
		expectedCount int
		hasError      bool
	}{
		{"split /24 to /26", "192.168.1.0/24", 26, 4, false},
		{"split /16 to /24", "10.0.0.0/16", 24, 256, false},
		{"invalid prefix", "192.168.1.0/24", 20, 0, true},
		{"IPv6 unsupported", "2001:db8::/64", 96, 0, true},
		{"invalid CIDR", "invalid", 26, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.cidr), lang.NumberValue(tt.newPrefixLen)}
			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			subnets := result.(lang.ListValue)
			if len(subnets) != tt.expectedCount {
				t.Errorf("Expected %d subnets, got %d", tt.expectedCount, len(subnets))
			}
		})
	}
}

func TestExpandIPv6(t *testing.T) {
	_, fn := expandIPv6()

	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"expand compressed IPv6", "2001:db8::1", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"expand loopback", "::1", "0000:0000:0000:0000:0000:0000:0000:0001", false},
		{"IPv4 address", "192.168.1.1", "", true},
		{"invalid IP", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestCompressIPv6(t *testing.T) {
	_, fn := compressIPv6()

	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{"compress expanded IPv6", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"already compressed", "2001:db8::1", false},
		{"IPv4 address", "192.168.1.1", true},
		{"invalid IP", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			compressed := string(result.(lang.StringValue))
			if len(compressed) == 0 {
				t.Errorf("Expected non-empty compressed IPv6")
			}
		})
	}
}

func TestIPToInt(t *testing.T) {
	_, fn := ipToInt()

	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"simple IPv4", "192.168.1.1", 3232235777, false},
		{"localhost", "127.0.0.1", 2130706433, false},
		{"zero IP", "0.0.0.0", 0, false},
		{"IPv6 address", "2001:db8::1", 0, true},
		{"invalid IP", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestIntToIP(t *testing.T) {
	_, fn := intToIP()

	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"simple conversion", 3232235777, "192.168.1.1"},
		{"localhost", 2130706433, "127.0.0.1"},
		{"zero", 0, "0.0.0.0"},
		{"max IPv4", 4294967295, "255.255.255.255"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestReverseIP(t *testing.T) {
	_, fn := reverseIP()

	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"IPv4 reverse", "192.168.1.1", "1.1.168.192.in-addr.arpa", false},
		{"IPv6 reverse", "2001:db8::1", "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa", false},
		{"invalid IP", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestIsRFC1918(t *testing.T) {
	_, fn := isRFC1918()

	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"10.x.x.x private", "10.0.0.1", true, false},
		{"172.16-31.x.x private", "172.16.0.1", true, false},
		{"192.168.x.x private", "192.168.1.1", true, false},
		{"public IP", "8.8.8.8", false, false},
		{"172.15.x.x not private", "172.15.0.1", false, false},
		{"172.32.x.x not private", "172.32.0.1", false, false},
		{"IPv6 address", "2001:db8::1", false, true},
		{"invalid IP", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestArgumentErrors(t *testing.T) {
	functions := []func() (string, lang.Function){
		isValidIP,
		cidrMatch,
		ipInRange,
	}

	for ti, tf := range functions {
		name, fn := tf()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			args := make([]lang.Value, ti)
			for i := range args {
				args[i] = lang.StringValue("test")
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"ipv4ToInt",
			func(t *testing.T) {
				ip := []byte{192, 168, 1, 1}
				result := ipv4ToInt(ip)
				expected := uint32(3232235777)
				if result != expected {
					t.Errorf("Expected %d, got %d", expected, result)
				}
			},
		},
		{
			"intToIPv4",
			func(t *testing.T) {
				result := intToIPv4(3232235777)
				expected := "192.168.1.1"
				if result.String() != expected {
					t.Errorf("Expected %s, got %s", expected, result.String())
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"isValidIP", "isIPv4", "isIPv6", "isPrivateIP", "isLoopbackIP",
		"isMulticastIP", "isLinkLocalIP", "cidrMatch", "cidrContains",
		"IPInRange", "cidrNetwork", "cidrBroadcast", "cidrHostCount",
		"cidrSubnets", "expandIPv6", "compressIPv6",
		"IPToInt", "intToIP", "reverseIP", "isRfc1918",
	}

	if len(functions) != len(expectedFunctions) {
		t.Errorf("Expected %d functions, got %d", len(expectedFunctions), len(functions))
	}

	for _, name := range expectedFunctions {
		if _, exists := functions[name]; !exists {
			t.Errorf("Expected function %s not found", name)
		}
	}
}

func BenchmarkIsValidIP(b *testing.B) {
	_, fn := isValidIP()
	args := []lang.Value{lang.StringValue("192.168.1.1")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkCidrMatch(b *testing.B) {
	_, fn := cidrMatch()
	args := []lang.Value{lang.StringValue("192.168.1.100"), lang.StringValue("192.168.1.0/24")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkIPToInt(b *testing.B) {
	_, fn := ipToInt()
	args := []lang.Value{lang.StringValue("192.168.1.1")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkIPInRange(b *testing.B) {
	_, fn := ipInRange()
	args := []lang.Value{
		lang.StringValue("192.168.1.100"),
		lang.StringValue("192.168.1.1"),
		lang.StringValue("192.168.1.200"),
	}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}
