package easyringbuffer

import (
	"errors"
	"sync"
)

// Error messages for ring buffer operations.
var (
	ErrRingBufferFull  = errors.New("ring buffer is full")
	ErrRingBufferEmpty = errors.New("ring buffer is empty")
)

// RingBuffer is a thread-safe ring buffer implementation.
type RingBuffer[T any] struct {
	buffer   []T
	capacity int
	size     int
	head     int
	tail     int
	mu       sync.Mutex
}

// New creates a new ring buffer with the specified capacity.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	return &RingBuffer[T]{
		buffer:   make([]T, capacity),
		capacity: capacity,
	}
}

// Enqueue adds an item to the ring buffer.
// Returns an error if the buffer is full.
func (rb *RingBuffer[T]) Enqueue(item T) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == rb.capacity {
		return ErrRingBufferFull
	}

	rb.buffer[rb.tail] = item
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size++

	return nil
}

// Dequeue removes and returns the oldest item from the ring buffer.
// Returns an error if the buffer is empty.
func (rb *RingBuffer[T]) Dequeue() (T, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	var zeroValue T
	if rb.size == 0 {
		return zeroValue, ErrRingBufferEmpty
	}

	item := rb.buffer[rb.head]
	// Optional: Clear the slot for garbage collection.
	rb.buffer[rb.head] = zeroValue
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--

	return item, nil
}

// GetAll returns all items in the buffer in order from oldest to newest.
func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	result := make([]T, rb.size)
	for i := 0; i < rb.size; i++ {
		index := (rb.head + i) % rb.capacity
		result[i] = rb.buffer[index]
	}
	return result
}

// GetLastN returns the last N items from the buffer.
func (rb *RingBuffer[T]) GetLastN(n int) []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if n > rb.size {
		n = rb.size
	}

	result := make([]T, n)
	for i := 0; i < n; i++ {
		index := (rb.tail - n + i + rb.capacity) % rb.capacity
		result[i] = rb.buffer[index]
	}
	return result
}

// Peek returns the next item without removing it from the buffer.
// Returns an error if the buffer is empty.
func (rb *RingBuffer[T]) Peek() (T, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	var zeroValue T
	if rb.size == 0 {
		return zeroValue, ErrRingBufferEmpty
	}

	return rb.buffer[rb.head], nil
}

// Len returns the current number of items in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.size
}

// Capacity returns the capacity of the ring buffer.
func (rb *RingBuffer[T]) Capacity() int {
	return rb.capacity
}

// IsEmpty checks if the ring buffer is empty.
func (rb *RingBuffer[T]) IsEmpty() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.size == 0
}

// IsFull checks if the ring buffer is full.
func (rb *RingBuffer[T]) IsFull() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.size == rb.capacity
}

// Reset clears all items in the ring buffer.
func (rb *RingBuffer[T]) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	var zeroValue T
	for i := 0; i < rb.capacity; i++ {
		rb.buffer[i] = zeroValue
	}
	rb.size = 0
	rb.head = 0
	rb.tail = 0
}
