# IP Package

The IP package provides comprehensive functions for IP address validation, manipulation, CIDR operations, and network calculations for both IPv4 and IPv6 addresses.

## Basic IP Validation

### `isValidIP(ip)`
Checks if a string is a valid IP address (IPv4 or IPv6).
- **Parameters:** `ip` (string) - IP address to validate
- **Returns:** Boolean indicating validity
- **Example:** `isValidIP("192.168.1.1")` → `true`

### `isIPv4(ip)`
Checks if a string is a valid IPv4 address.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's IPv4
- **Example:** `isIPv4("192.168.1.1")` → `true`

### `isIPv6(ip)`
Checks if a string is a valid IPv6 address.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's IPv6
- **Example:** `isIPv6("2001:db8::1")` → `true`

## IP Address Classification

### `isPrivateIP(ip)`
Checks if an IP address is in a private address range.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's private
- **Example:** `isPrivateIP("192.168.1.1")` → `true`

### `isLoopbackIP(ip)`
Checks if an IP address is a loopback address.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's loopback
- **Example:** `isLoopbackIP("127.0.0.1")` → `true`

### `isMulticastIP(ip)`
Checks if an IP address is a multicast address.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's multicast
- **Example:** `isMulticastIP("224.0.0.1")` → `true`

### `isLinkLocalIP(ip)`
Checks if an IP address is a link-local address.
- **Parameters:** `ip` (string) - IP address to check
- **Returns:** Boolean indicating if it's link-local
- **Example:** `isLinkLocalIP("169.254.1.1")` → `true`

### `isRfc1918(ip)`
Checks if an IPv4 address is in RFC 1918 private address space.
- **Parameters:** `ip` (string) - IPv4 address to check
- **Returns:** Boolean indicating if it's RFC 1918 private
- **Ranges:** 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
- **Example:** `isRfc1918("10.0.0.1")` → `true`

## CIDR Operations

### `cidrMatch(ip, cidr)`
Checks if an IP address falls within a CIDR range.
- **Parameters:** 
  - `ip` (string) - IP address to test
  - `cidr` (string) - CIDR notation (e.g., "192.168.1.0/24")
- **Returns:** Boolean indicating if IP is in CIDR range
- **Example:** `cidrMatch("192.168.1.100", "192.168.1.0/24")` → `true`

### `cidrContains(ip, cidr)`
Alias for `cidrMatch()` - checks if IP is contained in CIDR range.
- **Parameters:** Same as `cidrMatch()`
- **Returns:** Boolean indicating containment
- **Example:** `cidrContains("10.0.0.50", "10.0.0.0/16")` → `true`

### `cidrNetwork(cidr)`
Extracts the network address from a CIDR notation.
- **Parameters:** `cidr` (string) - CIDR notation
- **Returns:** Network address as string
- **Example:** `cidrNetwork("192.168.1.100/24")` → `"192.168.1.0"`

### `cidrBroadcast(cidr)`
Calculates the broadcast address for an IPv4 CIDR range.
- **Parameters:** `cidr` (string) - IPv4 CIDR notation
- **Returns:** Broadcast address as string
- **Note:** IPv6 doesn't have broadcast addresses
- **Example:** `cidrBroadcast("192.168.1.0/24")` → `"192.168.1.255"`

### `cidrHostCount(cidr)`
Calculates the number of host addresses in a CIDR range.
- **Parameters:** `cidr` (string) - CIDR notation
- **Returns:** Number of available host addresses
- **Note:** For IPv4, excludes network and broadcast addresses
- **Example:** `cidrHostCount("192.168.1.0/24")` → `254`

### `cidrSubnets(cidr, newPrefixLength)`
Splits a CIDR range into smaller subnets.
- **Parameters:** 
  - `cidr` (string) - Original CIDR range
  - `newPrefixLength` (number) - New prefix length for subnets
- **Returns:** Array of subnet CIDR strings
- **Note:** Currently supports IPv4 only
- **Example:** `cidrSubnets("192.168.1.0/24", 26)` → `["192.168.1.0/26", "192.168.1.64/26", "192.168.1.128/26", "192.168.1.192/26"]`

## IP Range Operations

### `IPInRange(ip, startIP, endIP)`
Checks if an IP address falls within a specified range.
- **Parameters:** 
  - `ip` (string) - IP address to test
  - `startIP` (string) - Start of IP range
  - `endIP` (string) - End of IP range
- **Returns:** Boolean indicating if IP is in range
- **Example:** `IPInRange("192.168.1.100", "192.168.1.1", "192.168.1.200")` → `true`

## IP Address Manipulation

### `expandIPv6(ip)`
Expands an IPv6 address to its full form.
- **Parameters:** `ip` (string) - IPv6 address to expand
- **Returns:** Fully expanded IPv6 address
- **Example:** `expandIPv6("2001:db8::1")` → `"2001:0db8:0000:0000:0000:0000:0000:0001"`

### `compressIPv6(ip)`
Compresses an IPv6 address to its shortest form.
- **Parameters:** `ip` (string) - IPv6 address to compress
- **Returns:** Compressed IPv6 address
- **Example:** `compressIPv6("2001:0db8:0000:0000:0000:0000:0000:0001")` → `"2001:db8::1"`

## IP Address Conversion

### `IPToInt(ip)`
Converts an IPv4 address to its integer representation.
- **Parameters:** `ip` (string) - IPv4 address
- **Returns:** Integer representation of IP
- **Example:** `IPToInt("192.168.1.1")` → `3232235777`

### `intToIP(integer)`
Converts an integer to its IPv4 address representation.
- **Parameters:** `integer` (number) - Integer to convert
- **Returns:** IPv4 address string
- **Example:** `intToIP(3232235777)` → `"192.168.1.1"`

### `reverseIP(ip)`
Generates the reverse DNS notation for an IP address.
- **Parameters:** `ip` (string) - IP address to reverse
- **Returns:** Reverse DNS notation string
- **Examples:** 
  - `reverseIP("192.168.1.1")` → `"1.1.168.192.in-addr.arpa"`
  - `reverseIP("2001:db8::1")` → `"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa"`

## Usage Notes

### CIDR Notation
- CIDR notation uses slash notation (e.g., "192.168.1.0/24")
- The number after the slash indicates the number of network bits
- Common IPv4 prefixes: /24 (256 addresses), /16 (65,536 addresses), /8 (16,777,216 addresses)

### IPv4 vs IPv6
- All functions automatically detect IP version unless specifically noted
- Some functions are IPv4-only (like `cidrBroadcast` and `cidrSubnets`)
- IPv6 addresses can be in compressed or expanded form

### Private Address Ranges
**IPv4 Private Ranges (RFC 1918):**
- 10.0.0.0/8 (10.0.0.0 - 10.255.255.255)
- 172.16.0.0/12 (172.16.0.0 - 172.31.255.255)
- 192.168.0.0/16 (192.168.0.0 - 192.168.255.255)

### Special Address Types
- **Loopback:** 127.0.0.0/8 (IPv4), ::1 (IPv6)
- **Link-Local:** 169.254.0.0/16 (IPv4), fe80::/10 (IPv6)
- **Multicast:** 224.0.0.0/4 (IPv4), ff00::/8 (IPv6)

### Error Handling
- Invalid IP addresses return appropriate error messages
- Functions return `false` for boolean checks on invalid inputs
- CIDR functions validate both IP and prefix length
- Range operations require IPs of the same version (IPv4 or IPv6)