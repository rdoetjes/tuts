package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
 * convertDotNotationToUInt32 converts an IP address from dot notation (e.g., "192.168.1.1") to a 32-bit unsigned integer using bitwise operations.
 *
 * Parameter:
 *   - ip: The IP address in dot notation.
 *
 * Returns:
 *   The IP address as a 32-bit unsigned integer.
 *
 * Calculation:
 *   - Each octet of the IP address is shifted left by a multiple of 8 bits and combined using bitwise OR.
 *     Example:
 *        IP:           192.168.1.1 => 11000000.10101000.00000001.00000001
 *        Converted:                => 3232235777 (in decimal)
 *
 * Process:
 *   - Split the IP into octets.
 *   - Convert each octet to an integer.
 *   - Shift and combine into a 32-bit integer.
 */
func convertDotNotationToUInt32(ip string) uint32 {
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

/*
 * u32UserToDotNotation converts a 32-bit unsigned integer into an IP address in dot notation using bitwise operations.
 *
 * Parameter:
 *   - ip: The IP address as a 32-bit unsigned integer.
 *
 * Returns:
 *   The IP address in dot notation.
 *
 * Calculation:
 *   - The integer is split into four 8-bit segments, each representing an octet of the IP address.
 *     Example:
 *        Integer:      3232235777 => 11000000.10101000.00000001.00000001
 *        Dot Notation:            => 192     .168     .1       .1
 *
 * Process:
 *   - Extract each octet using bit masking and shifting.
 *   - Format the octets into dot notation.
 */
func u32UserToDotNotation(ip uint32) string {
	octet1 := (ip >> 24) & 0xFF
	octet2 := (ip >> 16) & 0xFF
	octet3 := (ip >> 8) & 0xFF
	octet4 := ip & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", octet1, octet2, octet3, octet4)
}

/*
 * calculateNetworkAddress computes the network and broadcast addresses using bitwise operations.
 *
 * Parameters:
 *   - ipU32: The IP address as a 32-bit unsigned integer.
 *   - netmaskU32: The netmask as a 32-bit unsigned integer.
 *
 * Returns:
 *   - networkAddress: The network address as a 32-bit unsigned integer.
 *   - broadcastAddress: The broadcast address as a 32-bit unsigned integer.
 *
 * Calculation:
 *   - Network address is calculated using a bitwise AND between the IP address and netmask:
 *     Example:
 *        IP:           10.10.10.0   => 00001010.00001010.00001010.00000000
 *        Netmask:     255.255.254.0 => 11111111.11111111.11111110.00000000
 *        Network Addr:              => 00001010.00001010.00001010.00000000 (10.10.10.0)
 *
 *   - Broadcast address is calculated by inverting the netmask and performing a bitwise OR with the network address:
 *     Example:
 *        Network Addr:              => 00001010.00001010.00001010.00000000
 *        Inverted Netmask:          => 00000000.00000000.00000001.11111111
 *        Broadcast Addr:            => 00001010.00001010.00001011.11111111 (10.10.11.255)
 */
func calculateNetworkAddress(ipU32 uint32, netmaskU32 uint32) (uint32, uint32) {
	networkAddress := ipU32 & netmaskU32
	broadcastAddress := networkAddress | ^netmaskU32
	return networkAddress, broadcastAddress
}

/*
 * numberOfHosts calculates the number of usable hosts in a network using bitwise operations.
 *
 * Parameter:
 *   - netmaskU32: The netmask as a 32-bit unsigned integer.
 *
 * Returns:
 *   The total number of usable hosts as a 32-bit unsigned integer.
 *
 * Calculation:
 *   - The number of hosts is calculated by inverting the netmask to get the host bits and subtracting one:
 *     Example:
 *        Netmask:                255.255.255.0 => 11111111.11111111.11111111.00000000
 *        Inverted Netmask ^netmask:            => 00000000.00000000.00000000.11111111 (255)
 *        Usable Hosts:                         => 255 - 1 = 254
 */
func numberOfHosts(netmaskU32 uint32) uint32 {
	hostBits := ^netmaskU32
	return hostBits - 1
}

/*
 * create a 32bit variable
 * set all the bits to 1
 * 32-nr of set bits
 * shift that number of in (adding the 0s at the beginning)
 */
func convertCIDRToIPNetmask(cidr string) (string, string) {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return "", ""
	}

	ipU32 := convertDotNotationToUInt32(parts[0])
	prefixLength, err := strconv.Atoi(parts[1])
	if err != nil || prefixLength < 0 || prefixLength > 32 {
		fmt.Println("Invalid prefix length")
		return "", ""
	}

	var mask uint32 = ^uint32(0) << (32 - prefixLength)
	return u32UserToDotNotation(ipU32), u32UserToDotNotation(mask)
}

/*
 * Parse the CLI arguments make sure the input is correct, otherwise do not start
 */
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

/*
 * Just calls the calculation functions and prints the output
 */
func calcNetworkDetails(ip *string, netmask *string) {
	ipU32 := convertDotNotationToUInt32(*ip)
	netmaskU32 := convertDotNotationToUInt32(*netmask)
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
