package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	//calc "phonax.com/netcalc/lib"
)

func runCommand(arg ...string) ([]byte, error) {
	cmd := exec.Command("./bin/netcalc", arg...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("integration test failed to execute netcalc: %v, output: %s", err, string(output))
	}
	return output, nil
}

func TestNetcalcIntegrationCorrectIp(t *testing.T) {
	output, _ := runCommand("-ip", "192.168.10.139", "-netmask", "255.255.252.0")

	expectedOutput := "IP Address: 192.168.10.139 Netmask: 255.255.252.0\nNetwork Address: 192.168.8.0\nBroadcast Address: 192.168.11.255\nNumber of usable hosts: 1022\n"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Logf("Integration test passed")
	}
}

func TestNetcalcIntegrationCorrectCIDR(t *testing.T) {
	output, _ := runCommand("-cidr", "192.168.10.139/22")

	expectedOutput := "IP Address: 192.168.10.139 Netmask: 255.255.252.0\nNetwork Address: 192.168.8.0\nBroadcast Address: 192.168.11.255\nNumber of usable hosts: 1022\n"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Logf("Integration test passed")
	}
}

func TestNetcalcIntegrationWrongIp(t *testing.T) {
	output, _ := runCommand("-ip", "192.168.10.300", "-netmask", "255.255.252.0")

	expectedOutput := "IP Address incorrectly formatted, needs to be 4 octets 0-255"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Logf("Integration test passed")
	}
}

func TestNetcalcIntegrationWrongMask(t *testing.T) {
	output, _ := runCommand("-ip", "192.168.10.3", "-netmask", "255.255.252")

	expectedOutput := "Netmask incorrectly formatted, needs to be 4 octets 0-255"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Log("Integration test passed")
	}
}

func TestNetcalcIntegrationWrongCidr(t *testing.T) {
	output, _ := runCommand("-cidr", "192.168.10.3/64")
	expectedOutput := "CIDR mask incorrect should be 0-32"
	actualOutput := string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("1: Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Log("1: Integration test passed")
	}

	output, _ = runCommand("-cidr", "192.168.10.256/22")
	expectedOutput = "IP Address incorrectly formatted, needs to be 4 octets 0-255"
	actualOutput = string(output)

	if strings.TrimSpace(actualOutput) != strings.TrimSpace(expectedOutput) {
		t.Errorf("2: Integration test failed. Unexpected output: expected:\n%s\ngot:\n%s", expectedOutput, actualOutput)
	} else {
		t.Log("2: Integration test passed")
	}
}
