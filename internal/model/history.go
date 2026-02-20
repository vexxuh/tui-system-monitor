package model

// RingBuffer is an immutable ring buffer of float64 values.
// Push returns a new buffer; the original is never mutated.
type RingBuffer struct {
	data []float64
	size int
}

// NewRingBuffer creates a ring buffer with the given capacity.
func NewRingBuffer(size int) RingBuffer {
	if size <= 0 {
		size = 60
	}
	return RingBuffer{data: nil, size: size}
}

// Push returns a new RingBuffer with value appended.
// If at capacity, the oldest value is dropped.
func (rb RingBuffer) Push(value float64) RingBuffer {
	newLen := len(rb.data) + 1
	if newLen > rb.size {
		newLen = rb.size
	}
	newData := make([]float64, newLen)

	if len(rb.data) < rb.size {
		copy(newData, rb.data)
		newData[len(rb.data)] = value
	} else {
		copy(newData, rb.data[1:])
		newData[rb.size-1] = value
	}

	return RingBuffer{data: newData, size: rb.size}
}

// GetAll returns a copy of the buffer contents.
func (rb RingBuffer) GetAll() []float64 {
	if len(rb.data) == 0 {
		return nil
	}
	out := make([]float64, len(rb.data))
	copy(out, rb.data)
	return out
}

// Len returns the number of values currently in the buffer.
func (rb RingBuffer) Len() int {
	return len(rb.data)
}
