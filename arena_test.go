package uniquenumber

import (
	"testing"
	"unique"
)

type TestNumber interface {
	Signed | Unsigned | Floating
}

type TestHandle[T TestNumber] interface{ Value() T }

func runMake[Number TestNumber, Handle TestHandle[Number]](b *testing.B, makeHandle func(Number) Handle) {
	b.Helper()

	var cnt Number
	for b.Loop() {
		cnt++
		_ = makeHandle(cnt)
	}
}

func runValue[Number TestNumber, Handle TestHandle[Number]](b *testing.B, makeHandle func(Number) Handle) {
	b.Helper()

	h := makeHandle(67)
	for b.Loop() {
		_ = h.Value()
	}
}

func BenchmarkMake(b *testing.B) {
	b.ReportAllocs()

	b.Run("int", func(b *testing.B) { runMake(b, MakeInt[int64]) })
	b.Run("uint", func(b *testing.B) { runMake(b, MakeUint[uint64]) })
	b.Run("float", func(b *testing.B) { runMake(b, MakeFloat[float64]) })
	b.Run("unique", func(b *testing.B) { runMake(b, unique.Make[int]) })
	b.Run("new", func(b *testing.B) { runMake(b, makeNew[int]) })
}

func BenchmarkValue(b *testing.B) {
	b.ReportAllocs()

	b.Run("int", func(b *testing.B) { runValue(b, MakeInt[int64]) })
	b.Run("uint", func(b *testing.B) { runValue(b, MakeUint[uint64]) })
	b.Run("float", func(b *testing.B) { runValue(b, MakeFloat[float64]) })
	b.Run("unique", func(b *testing.B) { runValue(b, unique.Make[int]) })
	b.Run("new", func(b *testing.B) { runValue(b, makeNew[int]) })
}

type newVal[T any] struct {
	v T
}

func (n *newVal[T]) Value() T {
	return n.v
}

func makeNew[T any](v T) *newVal[T] {
	return &newVal[T]{v: v}
}

func TestMakeInt(t *testing.T) {
	const v int64 = 42
	h := MakeInt(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %d, got %d", v, got)
	}
}

func TestMakeIntZero(t *testing.T) {
	const v int64 = 0
	h := MakeInt(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %d, got %d", v, got)
	}
}

func TestMakeUint(t *testing.T) {
	const v uint64 = 1000
	h := MakeUint(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %d, got %d", v, got)
	}
}

func TestMakeUintZero(t *testing.T) {
	const v uint64 = 0
	h := MakeUint(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %d, got %d", v, got)
	}
}

func TestMakeFloatIntegral(t *testing.T) {
	const v float64 = 3.14
	h := MakeFloat(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %f, got %f", v, got)
	}
}

func TestMakeFloatZero(t *testing.T) {
	const v float64 = 0.0
	h := MakeFloat(v)
	if got := h.Value(); got != v {
		t.Errorf("expected %f, got %f", v, got)
	}
}
