//go:build freebsd

package main

import (
	"encoding/binary"
	"fmt"

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
	// Attachment logic for FreeBSD (requires PTRACE_ATTACH)
	return nil
}

func (s *FreeBSDScanner) Close() {}

func (s *FreeBSDScanner) InitialScan(target int32) ([]uintptr, error) {
	// FreeBSD memory scanning logic using /proc/<pid>/map and PTRACE_PEEKDATA
	return nil, fmt.Errorf("FreeBSD implementation pending")
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
