//go:build darwin

package main

/*
#include <mach/mach.h>
#include <mach/mach_vm.h>

// Wrapper for macros that cgo can't see directly
task_t get_self_task() {
    return mach_task_self();
}

// Wrapper to help with Mach calls
kern_return_t read_mem(task_t task, mach_vm_address_t addr, mach_vm_size_t size, vm_offset_t *data, mach_msg_type_number_t *cnt) {
    return mach_vm_read(task, addr, size, data, cnt);
}
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type DarwinScanner struct {
	task C.task_t
}

func NewPlatformScanner() MemoryScanner {
	return &DarwinScanner{}
}

func (s *DarwinScanner) Open(pid int) error {
	res := C.task_for_pid(C.get_self_task(), C.int(pid), &s.task)
	if res != C.KERN_SUCCESS {
		return fmt.Errorf("task_for_pid failed (code %d). Ensure you are running with sudo and have Developer Tools enabled", res)
	}
	return nil
}

func (s *DarwinScanner) Close() {}

func (s *DarwinScanner) InitialScan(target int32) ([]uintptr, error) {
	var results []uintptr
	var addr C.mach_vm_address_t = 0
	var size C.mach_vm_size_t

	for {
		var info C.vm_region_basic_info_data_64_t
		var count C.mach_msg_type_number_t = C.VM_REGION_BASIC_INFO_COUNT_64
		var objName C.mach_port_t

		res := C.mach_vm_region(s.task, &addr, &size, C.VM_REGION_BASIC_INFO_64, (C.vm_region_info_t)(unsafe.Pointer(&info)), &count, &objName)
		if res != C.KERN_SUCCESS {
			break
		}

		// Scan readable/writable memory
		if info.protection&C.VM_PROT_READ != 0 && info.protection&C.VM_PROT_WRITE != 0 {
			var data C.vm_offset_t
			var dataCnt C.mach_msg_type_number_t

			readRes := C.read_mem(s.task, addr, size, &data, &dataCnt)
			if readRes == C.KERN_SUCCESS {
				goData := C.GoBytes(unsafe.Pointer(uintptr(data)), C.int(dataCnt))
				for i := 0; i <= len(goData)-4; i += 4 {
					if int32(binary.LittleEndian.Uint32(goData[i:i+4])) == target {
						results = append(results, uintptr(addr)+uintptr(i))
					}
				}
				C.vm_deallocate(C.get_self_task(), data, C.vm_size_t(dataCnt))
			}
		}
		addr += size
	}
	return results, nil
}

func (s *DarwinScanner) Rescan(addresses []uintptr, target int32) ([]uintptr, error) {
	var results []uintptr
	for _, addr := range addresses {
		var data C.vm_offset_t
		var dataCnt C.mach_msg_type_number_t

		res := C.read_mem(s.task, C.mach_vm_address_t(addr), 4, &data, &dataCnt)
		if res == C.KERN_SUCCESS {
			val := int32(binary.LittleEndian.Uint32(C.GoBytes(unsafe.Pointer(uintptr(data)), 4)))
			if val == target {
				results = append(results, addr)
			}
			C.vm_deallocate(C.get_self_task(), data, C.vm_size_t(dataCnt))
		}
	}
	return results, nil
}
