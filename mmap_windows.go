//go:build windows

package uniquenumber

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func tryMmap(size int) unsafe.Pointer {
	addr, err := windows.VirtualAlloc(0, uintptr(size), windows.MEM_RESERVE, windows.PAGE_READONLY)
	if err != nil {
		return nil
	}
	return unsafe.Pointer(addr)
}
