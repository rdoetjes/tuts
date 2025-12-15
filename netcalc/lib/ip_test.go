package lib

import (
	"fmt"
	"testing"
	//calc "phonax.com/netcalc/lib"
)

func TestConvertIpToBinary(t *testing.T) {
	ip := "255.255.255.255"
	expected := uint32(4294967295)
	result := ConvertDotNotationToUInt32(ip)

	if result != expected {
		fmt.Printf("Test failed for ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		fmt.Println("Test passed for ConvertDotNotationToUInt32")
	}

	ip = "10.10.10.10"
	expected = uint32(168430090)
	result = ConvertDotNotationToUInt32(ip)

	if result != expected {
		fmt.Printf("Test failed for ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		fmt.Println("Test passed for ConvertDotNotationToUInt32")
	}
}

func TestConvertBinaryToIp(t *testing.T) {
	ipU32 := uint32(4294967295)
	expected := "255.255.255.255"
	result := U32UserToDotNotation(ipU32)

	if result != expected {
		fmt.Printf("Test failed for U32UserToDotNotation: expected %s, got %s\n", expected, result)
	} else {
		fmt.Println("Test passed for U32UserToDotNotation")
	}

	ipU32 = uint32(168430090)
	expected = "10.10.10.10"
	result = U32UserToDotNotation(ipU32)

	if result != expected {
		fmt.Printf("Test failed for U32UserToDotNotation: expected %s, got %s\n", expected, result)
	} else {
		fmt.Println("Test passed for U32UserToDotNotation")
	}
}

func TestConvCidrToIpNetmask(t *testing.T) {
	cidr := "192.168.178.0/24"
	expected_ip := "192.168.178.0"
	expected_netmask := "255.255.255.0"

	result_ip, result_netmask := ConvertCIDRToIPNetmask(cidr)

	if result_ip != expected_ip || result_netmask != expected_netmask {
		fmt.Printf("Test failed for ConvertCIDRToIPNetmask: expected %s %s, got %s %s\n", expected_ip, expected_netmask, result_ip, result_netmask)
	} else {
		fmt.Println("Test passed for ConvertCIDRToIPNetmask")
	}
}

func TestCalculateNetworkAddress(t *testing.T) {
	ipU32 := ConvertDotNotationToUInt32("192.168.10.139")
	netmaskU32 := ConvertDotNotationToUInt32("255.255.252.0")
	expected := ConvertDotNotationToUInt32("192.168.8.0")
	result := CalculateNetworkAddress(ipU32, netmaskU32)

	if result != expected {
		t.Errorf("Test failed for CalculateNetworkAddress: expected %s, got %s\n", U32UserToDotNotation(expected), U32UserToDotNotation(result))
	} else {
		fmt.Println("Test passed for CalculateNetworkAddress")
	}
}

func TestCalculateBroadcastAddress(t *testing.T) {
	ipU32 := ConvertDotNotationToUInt32("192.168.10.139")
	netmaskU32 := ConvertDotNotationToUInt32("255.255.252.0")
	expected := ConvertDotNotationToUInt32("192.168.11.255")
	result := CalculateBroadcastAddress(ipU32, netmaskU32)

	if result != expected {
		t.Errorf("Test failed for CalculateBroadcastAddress: expected %s, got %s\n", U32UserToDotNotation(expected), U32UserToDotNotation(result))
	} else {
		fmt.Println("Test passed for CalculateBroadcastAddress")
	}
}
