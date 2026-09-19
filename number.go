package uniquenumber

import (
	"math"
	"unique"
	"unsafe"
)

// ensures `handle` is the same size/layout as `unique.Handle[T]` (1 pointer)
// and is thus fine to do the unsafe cast
var (
	_ [unsafe.Sizeof(handle{})]struct{} = [unsafe.Sizeof(unique.Handle[any]{})]struct{}{}
)

// handle is an internally-allocated unique number handle. It is kept the same
// size and layout as unique.Handle[T] (a single pointer) so that it can be
// safely cast to and from unique.Handle[int64] / unique.Handle[uint64].
type handle struct {
	p unsafe.Pointer
}

// Signed is the constraint for all signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is the constraint for all unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Floating is the constraint for floating-point types.
type Floating interface {
	~float32 | ~float64
}

// Int is a unique handle for signed integer types.
//
// For values within the arena range, the handle pointer is computed directly
// from the value with no heap allocation; otherwise it falls back to a
// unique.Handle[int64].
type Int[T Signed] handle

// Value returns the underlying value of type T.
func (i Int[T]) Value() T {
	p := uintptr(i.p)
	if uintptr(minIntPtr) <= p && p <= uintptr(maxIntPtr) {
		return T(int64(p - uintptr(zeroIntPtr)))
	}
	return T(cast[unique.Handle[int64]](i.p).Value())
}

// MakeInt returns a unique Int handle for the given value v.
//
// If v is within the arena range, no allocation occurs and the returned
// handle is a pointer into the memory-mapped arena. Otherwise, it falls
// back to unique.Make(v).
func MakeInt[T Signed](v T) Int[T] {
	n := int64(v)
	if minInt <= n && n <= maxInt {
		return Int[T]{p: unsafe.Add(zeroIntPtr, n)}
	}
	return cast[Int[T]](unique.Make(n))
}

// Uint is a unique handle for unsigned integer types.
//
// For values within the arena range, the handle pointer is computed directly
// from the value with no heap allocation; otherwise it falls back to a
// unique.Handle[uint64].
type Uint[T Unsigned] handle

// Value returns the underlying value of type T.
func (u Uint[T]) Value() T {
	p := uintptr(u.p)
	if uintptr(minUintPtr) <= p && p <= uintptr(maxUintPtr) {
		return T(p - uintptr(minUintPtr))
	}
	return T(cast[unique.Handle[uint64]](u.p).Value())
}

// MakeUint returns a unique Uint handle for the given value v.
//
// If v is within the arena range, no allocation occurs and the returned
// handle is a pointer into the memory-mapped arena. Otherwise, it falls
// back to unique.Make(v).
func MakeUint[T Unsigned](v T) Uint[T] {
	n := uint64(v)
	if n <= maxUint {
		return Uint[T]{p: unsafe.Add(minUintPtr, n)}
	}
	return cast[Uint[T]](unique.Make(n))
}

// Float is a unique handle for floating-point types.
//
// Only integral values within the int64 range are stored in the arena; all
// other values (NaN, infinity, and non-integral floats) fall back to a
// unique.Handle[int64].
type Float[T Floating] handle

// Value returns the underlying value of type T.
func (f Float[T]) Value() T {
	p := uintptr(f.p)
	if uintptr(minIntPtr) <= p && p <= uintptr(maxIntPtr) {
		return T(int64(p - uintptr(zeroIntPtr)))
	}
	return T(cast[unique.Handle[float64]](f.p).Value())
}

// MakeFloat returns a unique Float handle for the given value v.
//
// Integral values within the int64 range are stored in the arena with
// no allocation; NaN, infinity, and non-integral values fall back to
// unique.Make(v).
func MakeFloat[T Floating](v T) Float[T] {
	n := float64(v)
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return cast[Float[T]](unique.Make(n))
	}
	if n < float64(minInt) || n > float64(maxInt) || math.Trunc(n) != n {
		return cast[Float[T]](unique.Make(n))
	}
	return Float[T]{p: unsafe.Add(zeroIntPtr, int64(n))}
}
