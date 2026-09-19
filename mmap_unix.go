//go:build unix

package uniquenumber

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func tryMmap(size int) unsafe.Pointer {
	s, err := unix.Mmap(-1, 0, size, unix.PROT_READ, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return nil
	}
	return unsafe.Pointer(unsafe.SliceData(s))
}
