package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	//calc "phonax.com/netcalc/lib"
)

func TestNetcalcIntegration(t *testing.T) {
	cmd := exec.Command("./bin/netcalc", "-ip", "192.168.10.139", "-netmask", "255.255.252.0")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Failed to execute netcalc: %v", err)
	}

	expectedOutput := "IP Address: 192.168.10.139 Netmask: 255.255.252.0\nNetwork Address: 192.168.8.0\nBroadcast Address: 192.168.11.255\nNumber of usable hosts: 1022\n"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		fmt.Printf("Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		fmt.Println("Integration test passed")
	}
}
