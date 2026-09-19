package uniquenumber

import (
	"unsafe"
)

var (
	arenaPtr, arenaBits = allocArena()

	maxUint    = uint64(1<<arenaBits - 1)
	minUintPtr = arenaPtr
	maxUintPtr = unsafe.Add(arenaPtr, 1<<arenaBits-1)

	minInt     = int64(-1 << (arenaBits - 1))
	maxInt     = int64(1<<(arenaBits-1) - 1)
	minIntPtr  = arenaPtr
	zeroIntPtr = unsafe.Add(minIntPtr, (1 << (arenaBits - 1)))
	maxIntPtr  = unsafe.Add(minIntPtr, (1<<arenaBits)-1)
)

func allocArena() (ptr unsafe.Pointer, bits int) {
	const minBits = 16 // 16k is cheap to allocate if mmap fails
	const maxBits = 46 // highest that usually succeeds on Linux
	for bits = maxBits; bits >= minBits; bits-- {
		ptr = tryMmap(1 << bits)
		if ptr != nil {
			return ptr, bits
		}
	}
	return unsafe.Pointer(unsafe.SliceData(make([]byte, 1<<minBits))), minBits
}

// cast casts v to type T
func cast[T, U any](v U) T {
	return *(*T)(unsafe.Pointer(&v))
}
