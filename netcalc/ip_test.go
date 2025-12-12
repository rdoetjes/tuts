package main

import (
	"fmt"
	"testing"
)

func TestConvertIpToBinary(t *testing.T) {
	ip := "255.255.255.255"
	expected := uint32(4294967295)
	result := converDotNotationToBinary(ip)

	if result != expected {
		fmt.Printf("Test failed: expected %d, got %d\n", expected, result)
	} else {
		fmt.Println("Test passed")
	}

	ip = "10.10.10.10"
	expected = uint32(168430090)
	result = converDotNotationToBinary(ip)

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
