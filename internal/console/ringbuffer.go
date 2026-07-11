package console

type RingBuffer[T any] struct {
	buffer []T
	size   int
	head   int
	full   bool
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	rb := RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
		head:   0,
		full:   false,
	}
	return &rb
}

func (rb *RingBuffer[T]) Push(item T) {
	rb.buffer[rb.head] = item

	// increment head value and wrap when equals size
	rb.head = (rb.head + 1) % rb.size
	if !rb.full && rb.head == 0 {
		rb.full = true
	}
}

func (rb *RingBuffer[T]) ToSlice() []T {
	var slice []T

	if !rb.full {
		slice = rb.buffer[:rb.head]
		return slice
	}

	slice = make([]T, rb.size)
	remainingIdx := rb.size - rb.head
	copy(slice, rb.buffer[rb.head:])                // copy everything after the head into first half
	copy(slice[remainingIdx:], rb.buffer[:rb.head]) // copy remaining before head into second half

	return slice
}
