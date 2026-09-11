package main

// Scanner defines the interface for memory operations on different OSes
type MemoryScanner interface {
	Open(pid int) error
	InitialScan(target int32) ([]uintptr, error)
	Rescan(addresses []uintptr, target int32) ([]uintptr, error)
	Close()
}

// Global scanner instance
var scanner MemoryScanner

func initScanner(pid int) error {
	// The specific implementation is chosen at compile time via build tags
	scanner = NewPlatformScanner()
	return scanner.Open(pid)
}
