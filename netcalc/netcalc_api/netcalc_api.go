package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	calc "phonax.com/netcalc/lib"
)

// ApiResponse structure for API responses.
type ApiResponse struct {
	IPAddress        string `json:"ip_address"`
	Netmask          string `json:"netmask"`
	NetworkAddress   string `json:"network_address"`
	BroadcastAddress string `json:"broadcast_address"`
	UsableHosts      uint32 `json:"usable_hosts"`
	Error            string `json:"error,omitempty"`
}

func checkValidity(w http.ResponseWriter, ip string, netmask string, cidr string) (string, string) {
	if (ip != "" && netmask != "" && cidr != "") || (cidr == "" && (ip == "" || netmask == "")) {
		http.Error(w, "Please provide either both ip & netmask or a cidr", http.StatusBadRequest)
		return "", ""
	}

	if cidr != "" {
		var err error
		ip, netmask, err = calc.ConvertCIDRToIPNetmask(cidr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return "", ""
		}
	}
	return ip, netmask
}

func convertInputToU32(w http.ResponseWriter, ip string, netmask string) (uint32, uint32) {
	ipU32, err := calc.ConvertDotNotationToUInt32(ip)
	if err != nil || ipU32 == 0 {
		http.Error(w, "IP Address incorrectly formatted, needs to be 4 octets 0-255", http.StatusBadRequest)
		return 0, 0
	}

	netmaskU32, err := calc.ConvertDotNotationToUInt32(netmask)
	if err != nil || netmaskU32 == 0 {
		http.Error(w, "Netmask incorrectly formatted, needs to be 4 octets 0-255", http.StatusBadRequest)
		return 0, 0
	}

	return ipU32, netmaskU32
}

// calculateHandler processes the API requests and calculates network details.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	netmask := r.URL.Query().Get("netmask")
	cidr := r.URL.Query().Get("cidr")

	ip, netmask = checkValidity(w, ip, netmask, cidr)
	if ip == "" || netmask == "" {
		return
	}

	ipU32, netmaskU32 := convertInputToU32(w, ip, netmask)
	if ipU32 == 0 || netmaskU32 == 0 {
		return
	}

	networkAddress := calc.CalculateNetworkAddress(ipU32, netmaskU32)
	broadcastAddress := calc.CalculateBroadcastAddress(ipU32, netmaskU32)
	usableHosts := calc.NumberOfHosts(netmaskU32)

	response := ApiResponse{
		IPAddress:        ip,
		Netmask:          netmask,
		NetworkAddress:   calc.U32UserToDotNotation(networkAddress),
		BroadcastAddress: calc.U32UserToDotNotation(broadcastAddress),
		UsableHosts:      usableHosts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// main creates and runs the HTTP server.
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/calculate", calculateHandler).Methods("GET")

	log.Println("Starting netcalc API on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
