package lib

import (
	"testing"
	//calc "phonax.com/netcalc/lib"
)

func TestConvertIpToBinary(t *testing.T) {
	ip := "255.255.255.255"
	expected := uint32(4294967295)
	result, err := ConvertDotNotationToUInt32(ip)

	if result != expected || err != nil {
		t.Errorf("Test 1: failed for ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		t.Log("Test 1: passed for ConvertDotNotationToUInt32")
	}

	ip = "10.10.10.10"
	expected = uint32(168430090)
	result, err = ConvertDotNotationToUInt32(ip)

	if result != expected || err != nil {
		t.Errorf("Test 2: failed for ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		t.Log("Test 2: passed ConvertDotNotationToUInt32")
	}

	ip = "1234.123.123.123"
	expected = 0
	result, err = ConvertDotNotationToUInt32(ip)
	if result != expected || err != nil {
		t.Errorf("Test 3: failed ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		t.Log("Test 3: passed ConvertDotNotationToUInt32")
	}

	ip = "123.123.123"
	expected = 0
	result, err = ConvertDotNotationToUInt32(ip)
	if result != expected || err != nil {
		t.Errorf("Test 4: failed ConvertDotNotationToUInt32: expected %d, got %d\n", expected, result)
	} else {
		t.Log("Test 4: passed ConvertDotNotationToUInt32")
	}
}

func TestConvertBinaryToIp(t *testing.T) {
	ipU32 := uint32(4294967295)
	expected := "255.255.255.255"
	result := U32UserToDotNotation(ipU32)

	if result != expected {
		t.Errorf("1: Test failed for U32UserToDotNotation: expected %s, got %s\n", expected, result)
	} else {
		t.Log("1: Test passed for U32UserToDotNotation")
	}

	ipU32 = uint32(168430090)
	expected = "10.10.10.10"
	result = U32UserToDotNotation(ipU32)

	if result != expected {
		t.Errorf("2: Test failed for U32UserToDotNotation: expected %s, got %s\n", expected, result)
	} else {
		t.Log("2: Test passed for U32UserToDotNotation")
	}
}

func TestConvCidrToIpNetmask(t *testing.T) {
	cidr := "192.168.178.0/24"
	expected_ip := "192.168.178.0"
	expected_netmask := "255.255.255.0"

	result_ip, result_netmask, err := ConvertCIDRToIPNetmask(cidr)

	if result_ip != expected_ip || result_netmask != expected_netmask || err != nil {
		t.Errorf("Test 1: failed for ConvertCIDRToIPNetmask: expected %s %s, got %s %s\n", expected_ip, expected_netmask, result_ip, result_netmask)
	} else {
		t.Log("Test 1: passed for ConvertCIDRToIPNetmask")
	}

	cidr = "192.168.1782.0/242"
	expected_ip = ""
	expected_netmask = ""

	result_ip, result_netmask, err = ConvertCIDRToIPNetmask(cidr)
	if err != nil {
		t.Errorf("Input 2: failed for ConvertCIDRToIPNetmask: %s", err)
		return
	}

	if result_ip != expected_ip || result_netmask != expected_netmask {
		t.Errorf("Input 2: failed for ConvertCIDRToIPNetmask: expected %s %s, got %s %s\n", expected_ip, expected_netmask, result_ip, result_netmask)
	} else {
		t.Log("Test 2: passed for ConvertCIDRToIPNetmask")
	}
}

func TestCalculateNetworkAddress(t *testing.T) {
	ipU32, err := ConvertDotNotationToUInt32("192.168.10.139")
	if err != nil {
		t.Errorf("Input 1: failed for ConvertDotNotationToUInt32: %s", err)
		return
	}

	netmaskU32, err := ConvertDotNotationToUInt32("255.255.252.0")
	if err != nil {
		t.Errorf("Input 1: failed for ConvertDotNotationToUInt32: %s", err)
		return
	}
	expected, err := ConvertDotNotationToUInt32("192.168.8.0")
	if err != nil {
		t.Errorf("Input 1: failed for ConvertDotNotationToUInt32: %s", err)
		return
	}
	result := CalculateNetworkAddress(ipU32, netmaskU32)

	if result != expected {
		t.Errorf("Test failed for CalculateNetworkAddress: expected %s, got %s\n", U32UserToDotNotation(expected), U32UserToDotNotation(result))
	} else {
		t.Log("Test passed for CalculateNetworkAddress")
	}
}

func TestCalculateBroadcastAddress(t *testing.T) {
	ipU32, _ := ConvertDotNotationToUInt32("192.168.10.139")
	netmaskU32, _ := ConvertDotNotationToUInt32("255.255.252.0")
	expected, _ := ConvertDotNotationToUInt32("192.168.11.255")
	result := CalculateBroadcastAddress(ipU32, netmaskU32)

	if result != expected {
		t.Errorf("Test failed for CalculateBroadcastAddress: expected %s, got %s\n", U32UserToDotNotation(expected), U32UserToDotNotation(result))
	} else {
		t.Log("Test passed for CalculateBroadcastAddress")
	}
}
