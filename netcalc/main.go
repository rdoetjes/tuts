package main

import (
	"flag"
	"fmt"
	"os"

	"phonax.com/netcalc/helpers"
)

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
		*ip, *netmask = helpers.ConvertCIDRToIPNetmask(*cidr)
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
	ipU32 := helpers.ConvertDotNotationToUInt32(*ip)
	netmaskU32 := helpers.ConvertDotNotationToUInt32(*netmask)
	networkAddress, broadcastAddress := helpers.CalculateNetworkAddress(ipU32, netmaskU32)

	fmt.Printf("IP Address: %s Netmask: %s\n", *ip, *netmask)
	fmt.Printf("Network Address: %s\n", helpers.U32UserToDotNotation(networkAddress))
	fmt.Printf("Broadcast Address: %s\n", helpers.U32UserToDotNotation(broadcastAddress))
	fmt.Printf("Number of usable hosts: %d\n", helpers.NumberOfHosts(netmaskU32))
}

func main() {
	var ip string
	var netmask string
	var cidr string

	parseCLI(&ip, &netmask, &cidr)
	calcNetworkDetails(&ip, &netmask)
}
