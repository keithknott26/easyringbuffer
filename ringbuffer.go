package easyringbuffer

import (
	"errors"
	"sync/atomic"
)

var (
	ErrRingBufferFull  = errors.New("ring buffer is full")
	ErrRingBufferEmpty = errors.New("ring buffer is empty")
)

// RingBuffer is a high-performance, thread-safe ring buffer.
type RingBuffer[T any] struct {
	buffer       []T
	mask         uint64
	writePointer uint64
	readPointer  uint64
	writeReserve uint64
	readReserve  uint64
}

// NewRingBuffer creates a new ring buffer with the specified capacity.
// The capacity must be a power of two for optimal performance.
func NewRingBuffer[T any](capacity uint64) (*RingBuffer[T], error) {
	if capacity == 0 || (capacity&(capacity-1)) != 0 {
		return nil, errors.New("capacity must be a power of two")
	}
	return &RingBuffer[T]{
		buffer: make([]T, capacity),
		mask:   capacity - 1,
	}, nil
}

// Enqueue adds an item to the ring buffer.
// Returns an error if the buffer is full.
func (rb *RingBuffer[T]) Enqueue(item T) error {
	for {
		wp := atomic.LoadUint64(&rb.writePointer)
		rp := atomic.LoadUint64(&rb.readReserve)
		if wp-rp >= uint64(len(rb.buffer)) {
			return ErrRingBufferFull
		}
		if atomic.CompareAndSwapUint64(&rb.writePointer, wp, wp+1) {
			index := wp & rb.mask
			rb.buffer[index] = item
			// Update the write reserve pointer
			for !atomic.CompareAndSwapUint64(&rb.writeReserve, wp, wp+1) {
				// Spin-wait
			}
			return nil
		}
		// Failed to reserve write pointer, retry
	}
}

// Dequeue removes and returns an item from the ring buffer.
// Returns an error if the buffer is empty.
func (rb *RingBuffer[T]) Dequeue() (T, error) {
	var zeroValue T
	for {
		rp := atomic.LoadUint64(&rb.readPointer)
		wp := atomic.LoadUint64(&rb.writeReserve)
		if rp >= wp {
			return zeroValue, ErrRingBufferEmpty
		}
		if atomic.CompareAndSwapUint64(&rb.readPointer, rp, rp+1) {
			index := rp & rb.mask
			item := rb.buffer[index]
			// Optional: Clear the slot
			var zero T
			rb.buffer[index] = zero
			// Update the read reserve pointer
			for !atomic.CompareAndSwapUint64(&rb.readReserve, rp, rp+1) {
				// Spin-wait
			}
			return item, nil
		}
		// Failed to reserve read pointer, retry
	}
}

// GetAt retrieves the item at the given index relative to the read pointer.
// Index 0 corresponds to the oldest item.
func (rb *RingBuffer[T]) GetAt(index int) (T, error) {
	var zeroValue T
	rp := atomic.LoadUint64(&rb.readReserve)
	wp := atomic.LoadUint64(&rb.writeReserve)
	size := wp - rp
	if uint64(index) >= size {
		return zeroValue, errors.New("index out of range")
	}
	actualIndex := (rp + uint64(index)) & rb.mask
	return rb.buffer[actualIndex], nil
}

// IsEmpty checks if the ring buffer is empty.
func (rb *RingBuffer[T]) IsEmpty() bool {
	rp := atomic.LoadUint64(&rb.readReserve)
	wp := atomic.LoadUint64(&rb.writePointer)
	return rp >= wp
}

// IsFull checks if the ring buffer is full.
func (rb *RingBuffer[T]) IsFull() bool {
	wp := atomic.LoadUint64(&rb.writePointer)
	rp := atomic.LoadUint64(&rb.readReserve)
	return wp-rp >= uint64(len(rb.buffer))
}

// Capacity returns the capacity of the ring buffer.
func (rb *RingBuffer[T]) Capacity() uint64 {
	return uint64(len(rb.buffer))
}

// Size returns the current number of items in the ring buffer.
func (rb *RingBuffer[T]) Size() uint64 {
	rp := atomic.LoadUint64(&rb.readReserve)
	wp := atomic.LoadUint64(&rb.writePointer)
	return wp - rp
}

// Reset clears all items in the ring buffer.
func (rb *RingBuffer[T]) Reset() {
	atomic.StoreUint64(&rb.readPointer, 0)
	atomic.StoreUint64(&rb.writePointer, 0)
	atomic.StoreUint64(&rb.readReserve, 0)
	atomic.StoreUint64(&rb.writeReserve, 0)
	// Clear buffer (optional)
	var zero T
	for i := range rb.buffer {
		rb.buffer[i] = zero
	}
}
