package main

import (
	"fmt"
	"testing"
)

func TestConvertIpToBinary(t *testing.T) {
	ip := "255.255.255.255"
	expected := uint32(4294967295)
	result := convertDotNotationToUInt32(ip)

	if result != expected {
		fmt.Printf("Test failed: expected %d, got %d\n", expected, result)
	} else {
		fmt.Println("Test passed")
	}

	ip = "10.10.10.10"
	expected = uint32(168430090)
	result = convertDotNotationToUInt32(ip)

	if result != expected {
		fmt.Printf("Test failed: expected %d, got %d\n", expected, result)
	} else {
		fmt.Println("Test passed")
	}
}

func TestConvertBinaryToIp(t *testing.T) {
	ipU32 := uint32(4294967295)
	expected := "255.255.255.255"
	result := u32UserToDotNotation(ipU32)

	if result != expected {
		fmt.Printf("Test failed: expected %s, got %s\n", expected, result)
	} else {
		fmt.Println("Test passed")
	}

	ipU32 = uint32(168430090)
	expected = "10.10.10.10"
	result = u32UserToDotNotation(ipU32)

	if result != expected {
		fmt.Printf("Test failed: expected %s, got %s\n", expected, result)
	} else {
		fmt.Println("Test passed")
	}
}

func TestConvCidrToIpNetmask(t *testing.T) {
	cidr := "192.168.178.0/24"
	expected_ip := "192.168.178.0"
	expected_netmask := "255.255.255.0"

	result_ip, result_netmask := convertCIDRToIPNetmask(cidr)

	if result_ip != expected_ip || result_netmask != expected_netmask {
		fmt.Printf("Test failed: expected %s %s, got %s %s\n", expected_ip, expected_netmask, result_ip, result_netmask)
	} else {
		fmt.Println("Test passed")
	}
}
