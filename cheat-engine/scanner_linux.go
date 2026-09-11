//go:build linux

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinuxScanner struct {
	pid int
}

func NewPlatformScanner() MemoryScanner {
	return &LinuxScanner{}
}

func (s *LinuxScanner) Open(pid int) error {
	s.pid = pid
	// Verify we can access the process memory
	_, err := os.Stat(fmt.Sprintf("/proc/%d/mem", pid))
	return err
}

func (s *LinuxScanner) Close() {}

func (s *LinuxScanner) InitialScan(target int32) ([]uintptr, error) {
	var results []uintptr

	// Open the maps file to find readable/writable regions
	mapsFile, err := os.Open(fmt.Sprintf("/proc/%d/maps", s.pid))
	if err != nil {
		return nil, err
	}
	defer mapsFile.Close()

	// Open the memory file for reading
	memFile, err := os.Open(fmt.Sprintf("/proc/%d/mem", s.pid))
	if err != nil {
		return nil, err
	}
	defer memFile.Close()

	scanner := bufio.NewScanner(mapsFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// Check permissions: must be readable and writable (rw-)
		perms := fields[1]
		if !strings.HasPrefix(perms, "rw") {
			continue
		}

		// Parse address range (e.g., 00400000-0040c000)
		addrRange := strings.Split(fields[0], "-")
		start, _ := strconv.ParseUint(addrRange[0], 16, 64)
		end, _ := strconv.ParseUint(addrRange[1], 16, 64)

		size := end - start
		if size <= 0 || size > 0x7FFFFFFF { // Sanity check on size
			continue
		}

		// Read the memory region
		data := make([]byte, size)
		_, err := memFile.ReadAt(data, int64(start))
		if err != nil {
			continue // Skip regions we can't read
		}

		// Scan for the 32-bit integer
		for i := 0; i <= len(data)-4; i += 4 {
			val := int32(binary.LittleEndian.Uint32(data[i : i+4]))
			if val == target {
				results = append(results, uintptr(start)+uintptr(i))
			}
		}
	}

	return results, nil
}

func (s *LinuxScanner) Rescan(addresses []uintptr, target int32) ([]uintptr, error) {
	memFile, err := os.Open(fmt.Sprintf("/proc/%d/mem", s.pid))
	if err != nil {
		return nil, err
	}
	defer memFile.Close()

	var results []uintptr
	data := make([]byte, 4)

	for _, addr := range addresses {
		_, err := memFile.ReadAt(data, int64(addr))
		if err != nil {
			continue
		}

		val := int32(binary.LittleEndian.Uint32(data))
		if val == target {
			results = append(results, addr)
		}
	}
	return results, nil
}
