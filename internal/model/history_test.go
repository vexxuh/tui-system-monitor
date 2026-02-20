package model

import (
	"testing"
)

func TestNewRingBufferEmpty(t *testing.T) {
	rb := NewRingBuffer(10)
	if rb.Len() != 0 {
		t.Errorf("expected len 0, got %d", rb.Len())
	}
	if got := rb.GetAll(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestPushImmutable(t *testing.T) {
	rb := NewRingBuffer(5)
	rb2 := rb.Push(1.0)
	rb3 := rb2.Push(2.0)

	if rb.Len() != 0 {
		t.Errorf("original mutated: len %d", rb.Len())
	}
	if rb2.Len() != 1 {
		t.Errorf("expected len 1, got %d", rb2.Len())
	}
	if rb3.Len() != 2 {
		t.Errorf("expected len 2, got %d", rb3.Len())
	}

	got := rb3.GetAll()
	if got[0] != 1.0 || got[1] != 2.0 {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestPushDropsOldest(t *testing.T) {
	rb := NewRingBuffer(3)
	for i := range 5 {
		rb = rb.Push(float64(i))
	}

	if rb.Len() != 3 {
		t.Errorf("expected len 3, got %d", rb.Len())
	}
	got := rb.GetAll()
	// Should have 2, 3, 4
	for i, want := range []float64{2, 3, 4} {
		if got[i] != want {
			t.Errorf("index %d: got %f, want %f", i, got[i], want)
		}
	}
}

func TestGetAllReturnsCopy(t *testing.T) {
	rb := NewRingBuffer(5).Push(1).Push(2).Push(3)
	got := rb.GetAll()
	got[0] = 999

	fresh := rb.GetAll()
	if fresh[0] != 1 {
		t.Error("GetAll did not return a copy")
	}
}

func TestDefaultSize(t *testing.T) {
	rb := NewRingBuffer(0)
	if rb.size != 60 {
		t.Errorf("expected default size 60, got %d", rb.size)
	}
}
