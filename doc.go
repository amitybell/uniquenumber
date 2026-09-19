// Package uniquenumber provides unique handles for signed integers, unsigned
// integers, and floating-point types that avoid heap allocation for values
// within a representable range.
//
// For values within the representable range, a handle is a compact pointer
// with no heap allocation. Values outside that range fall back to
// [unique.Make], which deduplicates via the runtime's unique map.
//
// This package is a drop-in replacement for [unique.Handle[T]] with lower
// allocation overhead for small integers.
//
// The minimum guaranteed representable range is 16-bit (±32,768 for signed,
// 0–65,535 for unsigned). The actual range may be larger on platforms that
// support it, up to 46-bit for unsigned and ±2^45 for signed.
package uniquenumber
