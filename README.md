# uniquenumber

`uniquenumber` is an int/uint/float-specific allocator that also provides unique handles for signed integers, unsigned integers, and floating-point numbers.

It can be used as a drop-in replacement for unique.Make/Handle when you need to intern numbers (the returned handle can be compared with `==`).

It's a fast zero-allocation allocator for numbers. The allocation strategy is essentially just projecting the number onto a 16-46-bit wide mmapped virtual address space (it doesn't actually allocate memory).

For numbers that fall within the representable range, it's zero-alloc. Otherwise, it falls back to `unique.Make`.

This was extracted from the faedra interpreter where it helped it reach speeds approaching Luajit in (non-jit) interpreter mode, so now your own programming language interpreter can finally be made fast as well :)

### Representable range

The minimum guaranteed representable range is 16-bit (±32,768 for signed, 0–65,535 for unsigned). The actual range is usually larger (if mmap succeeds) up to 46-bit for unsigned and ±2^45 for signed.

For floating-point types, only integral values within the above range are handled; NaN, infinity, and non-integral values fall back to `unique.Make`.

## Installation

```sh
go get github.com/amitybell/uniquenumber
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/amitybell/uniquenumber"
)

func main() {
	h1 := uniquenumber.MakeInt[int32](-67)
	fmt.Println(h1.Value()) // -67

	h2 := uniquenumber.MakeUint[uint32](67)
	fmt.Println(h2.Value()) // 1000

	h3 := uniquenumber.MakeFloat[float64](6.7)
	fmt.Println(h3.Value()) // 6.7

	h4 := uniquenumber.MakeInt[int64](1e12)
	fmt.Println(h4.Value()) // 1000000000000
}
```

## Comparison with `unique.Handle`

| Feature                   | `unique.Handle[T]` | `uniquenumber`                                  |
| ------------------------- | ------------------ | ----------------------------------------------- |
| Allocation (in-range)     | Heap allocation    | None (arena pointer)                            |
| Allocation (out-of-range) | Heap allocation    | Falls back to `unique.Make`                     |
| Types                     | Any `T`            | `~int*`, `~uint*`, `~float*` only               |
| Deduplication             | Runtime unique map | Arena for in-range, unique map for out-of-range |

## Benchmarks

The Int/Uint/Float numbers are usually ~2ns faster in real code (inlining becomes possible).
These numbers include the overhead due to using a benchmark helper function.

| Benchmark                            | Time/op     | Memory   | Allocations |
| ------------------------------------ | ----------- | -------- | ----------- |
| `Int/Make`                           | 4.2 ns/op   | 0 B/op   | 0 allocs/op |
| `Int/Value`                          | 2.2 ns/op   | 0 B/op   | 0 allocs/op |
| `Uint/Make`                          | 4.2 ns/op   | 0 B/op   | 0 allocs/op |
| `Uint/Value`                         | 2.2 ns/op   | 0 B/op   | 0 allocs/op |
| `Float/Make`                         | 6.5 ns/op   | 0 B/op   | 0 allocs/op |
| `Float/Value`                        | 2.4 ns/op   | 0 B/op   | 0 allocs/op |
| `Unique/Make` (`unique.Make[int64]`) | 1,272 ns/op | 180 B/op | 6 allocs/op |
| `Unique/Value`                       | 1.6 ns/op   | 0 B/op   | 0 allocs/op |
| `New/Make` (`new(int64)` reference)  | 10.9 ns/op  | 8 B/op   | 1 alloc/op  |
| `New/Value`                          | 1.6 ns/op   | 0 B/op   | 0 allocs/op |

`uniquenumber` achieves zero allocations for in-range values, compared to `unique.Make` which allocates 180 bytes and 6 objects per call.

## Disclaimer

I only tested it on Linux
