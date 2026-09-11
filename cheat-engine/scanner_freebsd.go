//go:build freebsd

package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"

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

// kinfo_vmentry structure for FreeBSD
type kinfoVmentry struct {
	StructSize   int32
	Type         int32
	Start        uint64
	End          uint64
	Offset       uint64
	VnodeId      uint64
	Status       uint32
	Extra3       uint32
	Read         int32
	Write        int32
	Execute      int32
	CopyOnWrite  int32
	NeedsCopy    int32
	TypeSpecific [112]byte // Simplified
}

func (s *FreeBSDScanner) InitialScan(target int32) ([]uintptr, error) {
	var results []uintptr

	// Use sysctl kern.proc.vmmap.<pid> to get memory maps without /proc
	// First call to get the required size
	bufSize, err := unix.SysctlRaw("kern.proc.vmmap", s.pid)
	if err != nil {
		return nil, fmt.Errorf("sysctl kern.proc.vmmap failed: %w", err)
	}

	if len(bufSize) == 0 {
		return nil, fmt.Errorf("no memory maps returned for pid %d", s.pid)
	}

	// The buffer contains a sequence of kinfo_vmentry structures
	// Each structure starts with its size (int32)
	offset := 0
	for offset < len(bufSize) {
		entry := (*kinfoVmentry)(unsafe.Pointer(&bufSize[offset]))

		if entry.StructSize == 0 {
			break
		}

		// Check permissions: must be readable and writable
		if entry.Read != 0 && entry.Write != 0 {
			start := uintptr(entry.Start)
			size := uintptr(entry.End - entry.Start)

			if size > 0 && size < 0x7FFFFFFF {
				regionResults := s.scanRegion(start, size, target)
				results = append(results, regionResults...)
			}
		}

		offset += int(entry.StructSize)
	}

	return results, nil
}

func (s *FreeBSDScanner) scanRegion(start uintptr, size uintptr, target int32) []uintptr {
	var results []uintptr

	// FreeBSD's PtracePeekData can be finicky with buffer sizes in some Go versions.
	// We'll use a larger buffer and read in chunks to be more efficient and avoid the slice panic.
	const chunkSize = 4096
	buf := make([]byte, chunkSize)

	for i := uintptr(0); i < size; i += chunkSize {
		remaining := size - i
		readSize := chunkSize
		if uintptr(readSize) > remaining {
			readSize = int(remaining)
		}

		// Ensure we read at least 4 bytes if possible
		if readSize < 4 {
			break
		}

		// Use PtraceIO directly for more control and to avoid the PtracePeekData wrapper issues
		piod := unix.PtraceIoDesc{
			Op:   unix.PT_IO_READ_D,
			Off:  uint64(start + i),
			Addr: uintptr(unsafe.Pointer(&buf[0])),
			Len:  uint(readSize),
		}

		err := unix.PtraceIO(s.pid, &piod)
		if err != nil {
			continue
		}

		// Scan the buffer we just read
		for j := 0; j <= int(piod.Len)-4; j += 4 {
			val := int32(binary.LittleEndian.Uint32(buf[j : j+4]))
			if val == target {
				results = append(results, start+i+uintptr(j))
			}
		}
	}
	return results
}

func (s *FreeBSDScanner) Rescan(addresses []uintptr, target int32) ([]uintptr, error) {
	var results []uintptr
	data := make([]byte, 4)

	for _, addr := range addresses {
		piod := unix.PtraceIoDesc{
			Op:   unix.PT_IO_READ_D,
			Off:  uint64(addr),
			Addr: uintptr(unsafe.Pointer(&data[0])),
			Len:  4,
		}

		err := unix.PtraceIO(s.pid, &piod)
		if err == nil && piod.Len == 4 {
			if int32(binary.LittleEndian.Uint32(data)) == target {
				results = append(results, addr)
			}
		}
	}
	return results, nil
}
