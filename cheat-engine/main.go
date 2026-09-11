package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cheat-engine <process_name>")
		os.Exit(1)
	}

	procName := os.Args[1]
	pid, err := getPidByName(procName)
	if err != nil {
		fmt.Printf("Error finding process: %v\n", err)
		os.Exit(1)
	}

	if err := initScanner(pid); err != nil {
		fmt.Printf("Error: %v. Try running with sudo.\n", err)
		os.Exit(1)
	}
	defer scanner.Close()

	fmt.Printf("Attached to %s (PID: %d)\n", procName, pid)

	reader := bufio.NewReader(os.Stdin)
	var matches []uintptr

	for {
		fmt.Printf("\n[%d matches] Enter value to scan ('q' to quit, 'r' to reset): ", len(matches))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			break
		}
		if input == "r" {
			matches = nil
			fmt.Println("Scanner reset.")
			continue
		}

		val, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			fmt.Println("Please enter a valid 32-bit integer.")
			continue
		}

		target := int32(val)

		if len(matches) == 0 {
			fmt.Println("Performing initial scan...")
			matches, err = scanner.InitialScan(target)
		} else {
			fmt.Printf("Rescanning %d potential addresses...\n", len(matches))
			matches, err = scanner.Rescan(matches, target)
		}

		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}

		fmt.Printf("Found %d matches.\n", len(matches))
		if len(matches) > 0 && len(matches) <= 20 {
			for _, addr := range matches {
				fmt.Printf("0x%x\n", addr)
			}
		}
	}
}

func getPidByName(name string) (int, error) {
	out, err := exec.Command("pgrep", "-x", name).Output()
	if err != nil {
		return 0, fmt.Errorf("process '%s' not found", name)
	}
	pidStr := strings.TrimSpace(string(out))
	return strconv.Atoi(pidStr)
}
