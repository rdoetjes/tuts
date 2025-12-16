package main

import (
	"flag"
	"fmt"
	"os"

	calc "phonax.com/netcalc/lib"
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
		if *ip, *netmask = calc.ConvertCIDRToIPNetmask(*cidr); *ip == "" || *netmask == "" {
			fmt.Println("Invalid CIDR")
			os.Exit(0)
		}
	}
}

func checkDotNotResult(s uint32, err string) {
	if s == 0 {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
 * Just calls the calculation functions and prints the output
 */
func calcNetworkDetails(ip *string, netmask *string) {
	ipU32 := calc.ConvertDotNotationToUInt32(*ip)
	checkDotNotResult(ipU32, "IP Address in correctly formatted, needs to be 4 octets 0-255")

	netmaskU32 := calc.ConvertDotNotationToUInt32(*netmask)
	checkDotNotResult(netmaskU32, "Netmask in correctly formatted")

	networkAddress := calc.CalculateNetworkAddress(ipU32, netmaskU32)
	broadcastAddress := calc.CalculateBroadcastAddress(ipU32, netmaskU32)

	fmt.Printf("IP Address: %s Netmask: %s\n", *ip, *netmask)
	fmt.Printf("Network Address: %s\n", calc.U32UserToDotNotation(networkAddress))
	fmt.Printf("Broadcast Address: %s\n", calc.U32UserToDotNotation(broadcastAddress))
	fmt.Printf("Number of usable hosts: %d\n", calc.NumberOfHosts(netmaskU32))
}

func main() {
	var ip string
	var netmask string
	var cidr string

	parseCLI(&ip, &netmask, &cidr)
	calcNetworkDetails(&ip, &netmask)
}
