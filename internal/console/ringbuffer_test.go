package console

import (
	"testing"
)

func TestRingBuffer_UnfilledToSlice(t *testing.T) {
	rb := NewRingBuffer[string](5)
	rb.Push("a")
	rb.Push("b")

	got := rb.ToSlice()
	want := []string{"a", "b"}

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestRingBuffer_WrapAround(t *testing.T) {
	rb := NewRingBuffer[string](3)
	rb.Push("a")
	rb.Push("b")
	rb.Push("c")
	rb.Push("d") // overwrites "a"

	got := rb.ToSlice()
	want := []string{"b", "c", "d"}

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestRingBuffer_FullBuffer(t *testing.T) {
	rb := NewRingBuffer[string](3)
	rb.Push("a")
	rb.Push("b")
	rb.Push("c")

	if !rb.full {
		t.Error("expected buffer to be marked full")
	}
}

func TestRingBuffer_Empty(t *testing.T) {
	rb := NewRingBuffer[string](3)
	got := rb.ToSlice()
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}
