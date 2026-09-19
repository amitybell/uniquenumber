//go:build !unix && !windows

package uniquenumber

import (
	"unsafe"
)

func tryMmap(_ int) unsafe.Pointer {
	return nil
}
