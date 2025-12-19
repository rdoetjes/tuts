package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// areOctetsValuesCorrect checks if the values in an IPv4 octet are between 0 and 255.
// It returns true if they are, or false if a value is smaller than 0 or larger than 255.
// It is a private helper function used in the exported functions.
func areOctetsValuesCorrect(s string) (bool, error) {
	sOctets := strings.Split(s, ".")
	for _, sOctet := range sOctets {
		x, err := strconv.Atoi(sOctet)
		if err != nil || !isValueInRange(x, 0, 255) {
			return false, err
		}
	}
	return true, nil
}

// isCIDRMaskCorrect checks if the string value for the CIDR mask is between 0 and 32.
// It returns true if it is, or false if a value is smaller than 0 or larger than 32.
// It is a private helper function used in the exported functions.
func isCIDRMaskCorrect(s string) (bool, error) {
	mask, err := strconv.Atoi(s)
	if err != nil {
		return false, err
	}
	return isValueInRange(mask, 0, 32), nil
}

// isValueInRange checks if a value is within a specified range.
func isValueInRange(value, min, max int) bool {
	return value >= min && value <= max
}

// isValidIpOrNetMask checks whether a given string represents a valid IPv4 octet.
//
// The function evaluates if the provided string `sOctet` is a string that consists of
// 4 octets between 0-255 so 0-255.0-255.0-255.0-255
//
// Parameters:
//   - sOctet: A string representing the 4 octets value to be checked.
//
// Returns:
//   - A boolean: True if the `sOctet` is a IPv4 address or subnetmask, otherwise false.
//
// Example:
//
//	isValid := isValidIpOrNetMask("192.12.14.15")
//	isValid will be true because is a valid IPv4 octet.
//
//	isValid := isValidIpOrNetMask("192.12.1422.15")
//	isValid will be false because an octet must lay between 0-255.
func isValidIpOrNetMask(sOctet string) (bool, error) {
	regex := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`

	match, err := regexp.MatchString(regex, sOctet)
	if err != nil {
		return false, err
	}

	success, err := areOctetsValuesCorrect(sOctet)
	if err != nil || !success {
		return false, err
	}

	return match, nil
}

// ConvertDotNotationToUInt32 converts an IP address from dot notation (e.g., "192.168.1.1") to a 32-bit unsigned integer using bitwise operations.
// Parameter:
//   - ip: The IP address in dot notation.
//
// Returns:
//
//	The IP address as a 32-bit unsigned integer.
func ConvertDotNotationToUInt32(ip string) (uint32, error) {
	if success, err := isValidIpOrNetMask(ip); err != nil || !success {
		return 0, err
	}

	octets := strings.Split(ip, ".")
	var u32Ip uint32 = 0
	for _, octet := range octets {
		num, err := strconv.Atoi(octet)
		if err != nil {
			return 0, fmt.Errorf("invalid octet %s: %v", octet, err)
		}
		u32Ip = (u32Ip << 8) | uint32(num)
	}

	return u32Ip, nil
}

// U32UserToDotNotation converts a 32-bit unsigned integer into an IP address in dot notation using bitwise operations.
// Parameter:
//   - ip: The IP address as a 32-bit unsigned integer.
//
// Returns:
//
//	The IP address in dot notation.
func U32UserToDotNotation(ip uint32) string {
	octet1 := (ip >> 24) & 0xFF
	octet2 := (ip >> 16) & 0xFF
	octet3 := (ip >> 8) & 0xFF
	octet4 := ip & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", octet1, octet2, octet3, octet4)
}

// CalculateNetworkAddress computes the network address using bitwise operations.
// Parameters:
//   - ipU32: The IP address as a 32-bit unsigned integer.
//   - netmaskU32: The netmask as a 32-bit unsigned integer.
//
// Returns:
//
//	The network address as a 32-bit unsigned integer.
func CalculateNetworkAddress(ipU32 uint32, netmaskU32 uint32) uint32 {
	return ipU32 & netmaskU32
}

// CalculateBroadcastAddress computes the broadcast address using bitwise operations.
// Parameters:
//   - ipU32: The IP address as a 32-bit unsigned integer.
//   - netmaskU32: The netmask as a 32-bit unsigned integer.
//
// Returns:
//
//	The broadcast address as a 32-bit unsigned integer.
func CalculateBroadcastAddress(ipU32 uint32, netmaskU32 uint32) uint32 {
	networkAddress := CalculateNetworkAddress(ipU32, netmaskU32)
	return networkAddress | ^netmaskU32
}

// NumberOfHosts calculates the number of usable hosts in a network using bitwise operations.
// Parameter:
//   - netmaskU32: The netmask as a 32-bit unsigned integer.
//
// Returns:
//
//	The total number of usable hosts as a 32-bit unsigned integer.
func NumberOfHosts(netmaskU32 uint32) uint32 {
	hostBits := ^netmaskU32
	return hostBits - 1
}

// isValidCIDR checks whether a given string represents a valid CIDR notation.
//
// The function assesses if the `cidr` string is a valid representation of CIDR
// notation, consisting of a valid IPv4 address followed by a slash and a valid
// prefix length (0-32). It uses a regex pattern to perform this validation.
//
// Parameters:
//   - cidr: A string representing the CIDR notation to be validated.
//
// Returns:
//   - A boolean: True if the provided `cidr` is a valid CIDR notation; otherwise, false.
//   - an error: returns in case of regexp errors or Atoi conversion errors bubbling up from checks
//
// The function verifies:
//   - The IPv4 address part consists of four octets, each ranging from 0 to 255.
//   - The CIDR prefix length is a number between 0 and 32 inclusive.
//
// Example:
//
//	isValid := isValidCIDR("192.168.1.0/24")
//	isValid will be true because "192.168.1.0/24" is a valid CIDR notation.
//
//	isValid := isValidCIDR("256.100.100.0/24")
//	isValid will be false because the octet 256 is out of range.
func isValidCIDR(cidr string) (bool, error) {
	regex := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}$`
	match, err := regexp.MatchString(regex, cidr)

	if err != nil {
		return false, err
	}

	if !match {
		return false, fmt.Errorf("Invalid CIDR format 192.168.178.12/24")
	}

	parts := strings.Split(cidr, "/")

	if success, err := areOctetsValuesCorrect(parts[0]); err != nil || !success {
		return false, fmt.Errorf("IP Address incorrectly formatted, needs to be 4 octets 0-255")
	}

	if success, err := isCIDRMaskCorrect(parts[1]); err != nil || !success {
		return false, fmt.Errorf("CIDR mask incorrect should be 0-32")
	}

	return match, nil
}

// ConvertCIDRToIPNetmask converts a CIDR notation string to an IP address and its corresponding netmask in dot-decimal notation.
// Parameter:
//   - cidr: The CIDR notation string (e.g., "192.168.1.1/24").
//
// Returns:
//   - ip: The IP address in dot-decimal notation.
//   - netmask: The netmask in dot-decimal notation.
//   - error: Can be conversion errors bubbling up from checks or IP/Netmask format error.
func ConvertCIDRToIPNetmask(cidr string) (string, string, error) {
	if succcess, err := isValidCIDR(cidr); err != nil || !succcess {
		return "", "", err
	}

	parts := strings.Split(cidr, "/")

	ipU32, err := ConvertDotNotationToUInt32(parts[0])
	if err != nil {
		return "", "", err
	}

	prefixLength, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("conversion error %v", err)
	}

	if prefixLength < 0 || prefixLength > 32 {
		return "", "", fmt.Errorf("invalid CIDR mask length")
	}

	var mask uint32 = ^uint32(0) << (32 - prefixLength)
	return U32UserToDotNotation(ipU32), U32UserToDotNotation(mask), nil
}
