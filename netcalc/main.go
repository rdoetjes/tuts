package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func converDotNotationToBinary(ip string) uint32 {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return 0
	}

	var u32Ip uint32 = 0
	for _, octet := range octets {
		num, _ := strconv.Atoi(octet)
		u32Ip = (u32Ip << 8) | uint32(num)
	}

	return u32Ip
}

func u32UserToDotNotation(ip uint32) string {
	octet1 := (ip >> 24) & 0xFF
	octet2 := (ip >> 16) & 0xFF
	octet3 := (ip >> 8) & 0xFF
	octet4 := ip & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", octet1, octet2, octet3, octet4)
}

func calculateNetworkAddress(ipU32 uint32, netmaskU32 uint32) (uint32, uint32) {
	networkAddress := ipU32 & netmaskU32
	broadcastAddress := networkAddress | ^netmaskU32
	return networkAddress, broadcastAddress
}

func numberOfHosts(netmaskU32 uint32) uint32 {
	hostBits := ^netmaskU32
	return hostBits - 1
}

func convertCIDRToIPNetmask(cidr string) (string, string) {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return "", ""
	}

	ipU32 := converDotNotationToBinary(parts[0])
	prefixLength, err := strconv.Atoi(parts[1])
	if err != nil || prefixLength < 0 || prefixLength > 32 {
		fmt.Println("Invalid prefix length")
		return "", ""
	}

	var mask uint32 = ^uint32(0) << (32 - prefixLength)
	return u32UserToDotNotation(ipU32), u32UserToDotNotation(mask)
}

func parseCLI(ip *string, netmask *string, cidr *string) {
	flag.StringVar(ip, "ip", "", "IP address in dot notation (e.g., 192.168.1.1)")
	flag.StringVar(netmask, "netmask", "", "Netmask in dot notation (e.g., 255.255.255.0)")
	flag.StringVar(cidr, "cidr", "", "Network in CIDR notation (e.g., 192.168.1.0/24)")

	flag.Parse()

	if (*ip != "" && *netmask != "" && *cidr != "") || (*cidr == "" && (*ip == "" || *netmask == "")) {
		fmt.Println("Please provide either both -ip <ip address> -netmask <netmask> or a -cidr <cidr>")
		os.Exit(0)
	}

	if *cidr != "" {
		*ip, *netmask = convertCIDRToIPNetmask(*cidr)
		if *ip == "" || *netmask == "" {
			fmt.Println("Invalid CIDR or ip and netmask")
			os.Exit(0)
		}
	}
}

func calcNetworkDetails(ip *string, netmask *string) {
	ipU32 := converDotNotationToBinary(*ip)
	netmaskU32 := converDotNotationToBinary(*netmask)
	networkAddress, broadcastAddress := calculateNetworkAddress(ipU32, netmaskU32)

	fmt.Printf("IP Address: %s Netmask: %s\n", *ip, *netmask)
	fmt.Printf("Network Address: %s\n", u32UserToDotNotation(networkAddress))
	fmt.Printf("Broadcast Address: %s\n", u32UserToDotNotation(broadcastAddress))
	fmt.Printf("Number of usable hosts: %d\n", numberOfHosts(netmaskU32))
}

func main() {
	var ip string
	var netmask string
	var cidr string

	parseCLI(&ip, &netmask, &cidr)
	calcNetworkDetails(&ip, &netmask)
}
