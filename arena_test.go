package uniquenumber

import (
	"testing"
	"unique"
)

func handleBench[Number Signed | Unsigned | Floating, Handle interface{ Value() Number }](b *testing.B, makeHande func(Number) Handle) {
	b.Helper()
	b.ReportAllocs()

	b.Run("Make", func(b *testing.B) {
		var cnt Number
		var res any
		for b.Loop() {
			cnt++
			res = makeHande(cnt)
		}
		h := res.(Handle)
		if r := h.Value(); r != cnt {
			b.Fatalf("expected: %v; got %v", cnt, r)
		}
	})

	b.Run("Value", func(b *testing.B) {
		h := MakeInt(67)
		for b.Loop() {
			_ = h.Value()
		}
	})
}

func BenchmarkInt(b *testing.B) {
	handleBench(b, MakeInt[int64])
}

func BenchmarkUint(b *testing.B) {
	handleBench(b, MakeUint[uint64])
}

func BenchmarkFloat(b *testing.B) {
	handleBench(b, MakeFloat[float64])
}

func BenchmarkUnique(b *testing.B) {
	handleBench(b, unique.Make[int64])
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

func BenchmarkNew(b *testing.B) {
	// just for reference in-case unique'ness isn't a requirement
	handleBench(b, makeNew[int64])
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
