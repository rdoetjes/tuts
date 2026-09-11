//go:build freebsd

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type FreeBSDScanner struct {
	pid int
}

func NewPlatformScanner() MemoryScanner {
	return &FreeBSDScanner{}
}

func (s *FreeBSDScanner) Open(pid int) error {
	s.pid = pid
	// Attach to the process to allow memory access
	err := unix.PtraceAttach(pid)
	if err != nil {
		return fmt.Errorf("failed to attach to process %d: %w", pid, err)
	}

	// Wait for the process to stop
	var status unix.WaitStatus
	_, err = unix.Wait4(pid, &status, 0, nil)
	if err != nil {
		unix.PtraceDetach(pid)
		return fmt.Errorf("failed to wait for process: %w", err)
	}

	return nil
}

func (s *FreeBSDScanner) Close() {
	if s.pid != 0 {
		unix.PtraceDetach(s.pid)
	}
}

func (s *FreeBSDScanner) InitialScan(target int32) ([]uintptr, error) {
	var results []uintptr

	// FreeBSD procfs provides memory maps in /proc/<pid>/map
	// Format: start end resident priv_resident obj_id perms ref_cnt shadow_cnt flags type
	mapPath := fmt.Sprintf("/proc/%d/map", s.pid)
	mapFile, err := os.Open(mapPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", mapPath, err)
	}
	defer mapFile.Close()

	scanner := bufio.NewScanner(mapFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// fields[5] contains permissions like "rw-"
		perms := fields[5]
		if !strings.Contains(perms, "r") || !strings.Contains(perms, "w") {
			continue
		}

		// fields[0] is start, fields[1] is end (hex with 0x prefix)
		start, err := strconv.ParseUint(strings.TrimPrefix(fields[0], "0x"), 16, 64)
		if err != nil {
			continue
		}
		end, err := strconv.ParseUint(strings.TrimPrefix(fields[1], "0x"), 16, 64)
		if err != nil {
			continue
		}

		size := end - start
		if size <= 0 || size > 0x7FFFFFFF {
			continue
		}

		// Read region using PtracePeekData in 4-byte chunks
		// Note: Using /proc/<pid>/mem would be faster if enabled (procfs is often restricted)
		// For robustness on FreeBSD, we'll try to use /proc/<pid>/mem first, fallback to Ptrace
		regionResults := s.scanRegion(uintptr(start), uintptr(size), target)
		results = append(results, regionResults...)
	}

	return results, nil
}

func (s *FreeBSDScanner) scanRegion(start uintptr, size uintptr, target int32) []uintptr {
	var results []uintptr
	data := make([]byte, 4)

	// Try reading /proc/<pid>/mem if available
	memPath := fmt.Sprintf("/proc/%d/mem", s.pid)
	memFile, err := os.Open(memPath)
	if err == nil {
		defer memFile.Close()
		buffer := make([]byte, size)
		_, err = memFile.ReadAt(buffer, int64(start))
		if err == nil {
			for i := 0; i <= len(buffer)-4; i += 4 {
				val := int32(binary.LittleEndian.Uint32(buffer[i : i+4]))
				if val == target {
					results = append(results, start+uintptr(i))
				}
			}
			return results
		}
	}

	// Fallback to PtracePeekData (much slower but reliable)
	for i := uintptr(0); i <= size-4; i += 4 {
		_, err := unix.PtracePeekData(s.pid, start+i, data)
		if err != nil {
			continue
		}
		val := int32(binary.LittleEndian.Uint32(data))
		if val == target {
			results = append(results, start+i)
		}
	}
	return results
}

func (s *FreeBSDScanner) Rescan(addresses []uintptr, target int32) ([]uintptr, error) {
	var results []uintptr
	for _, addr := range addresses {
		data := make([]byte, 4)
		_, err := unix.PtracePeekData(s.pid, addr, data)
		if err == nil {
			if int32(binary.LittleEndian.Uint32(data)) == target {
				results = append(results, addr)
			}
		}
	}
	return results, nil
}
