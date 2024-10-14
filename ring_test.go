package easyringbuffer

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	capacity := uint64(8)
	rb, err := NewRingBuffer[int](capacity)
	if err != nil {
		t.Fatalf("Failed to create ring buffer: %v", err)
	}

	// Test Enqueue and Dequeue
	for i := 0; i < int(capacity); i++ {
		if err := rb.Enqueue(i); err != nil {
			t.Fatalf("Failed to enqueue: %v", err)
		}
	}

	if !rb.IsFull() {
		t.Fatalf("Ring buffer should be full")
	}

	for i := 0; i < int(capacity); i++ {
		item, err := rb.Dequeue()
		if err != nil {
			t.Fatalf("Failed to dequeue: %v", err)
		}
		if item != i {
			t.Fatalf("Expected %d, got %d", i, item)
		}
	}

	if !rb.IsEmpty() {
		t.Fatalf("Ring buffer should be empty")
	}
}

func TestRingBufferConcurrent(t *testing.T) {
	capacity := uint64(1024)
	rb, err := NewRingBuffer[int](capacity)
	if err != nil {
		t.Fatalf("Failed to create ring buffer: %v", err)
	}

	var wg sync.WaitGroup
	numProducers := 4
	numConsumers := 4
	itemsPerProducer := 10000

	// Producers
	for p := 0; p < numProducers; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				for {
					if err := rb.Enqueue(p*itemsPerProducer + i); err == nil {
						break
					}
					// Buffer is full, retry
				}
			}
		}(p)
	}

	// Consumers
	consumedItems := make([]int, numProducers*itemsPerProducer)
	var consumedCount int64
	for c := 0; c < numConsumers; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				item, err := rb.Dequeue()
				if err != nil {
					// Buffer is empty, check if producers are done
					if atomic.LoadInt64(&consumedCount) >= int64(numProducers*itemsPerProducer) {
						return
					}
					continue
				}
				index := atomic.AddInt64(&consumedCount, 1) - 1
				consumedItems[index] = item
			}
		}()
	}

	wg.Wait()

	// Verify all items were consumed
	if int(consumedCount) != numProducers*itemsPerProducer {
		t.Fatalf("Expected %d items consumed, got %d", numProducers*itemsPerProducer, consumedCount)
	}
}

func TestRingBufferFull(t *testing.T) {
	capacity := uint64(8)
	rb, err := NewRingBuffer[int](capacity)
	if err != nil {
		t.Fatalf("Failed to create ring buffer: %v", err)
	}

	// Enqueue items until the buffer is full
	for i := 0; i < int(capacity); i++ {
		if err := rb.Enqueue(i); err != nil {
			t.Fatalf("Failed to enqueue item %d: %v", i, err)
		}
	}

	// Buffer should be full now
	if !rb.IsFull() {
		t.Fatalf("Ring buffer should be full")
	}

	// Attempt to enqueue another item, should fail
	err = rb.Enqueue(100)
	if err == nil {
		t.Fatalf("Expected error when enqueueing to full buffer, got nil")
	}
	if !errors.Is(err, ErrRingBufferFull) {
		t.Fatalf("Expected ErrRingBufferFull, got %v", err)
	}

	// Dequeue one item
	item, err := rb.Dequeue()
	if err != nil {
		t.Fatalf("Failed to dequeue item: %v", err)
	}
	if item != 0 {
		t.Fatalf("Expected item 0, got %d", item)
	}

	// Buffer should not be full now
	if rb.IsFull() {
		t.Fatalf("Ring buffer should not be full after dequeueing an item")
	}

	// Enqueue another item, should succeed
	if err := rb.Enqueue(200); err != nil {
		t.Fatalf("Failed to enqueue item after dequeueing: %v", err)
	}
}

func TestRingBufferEmpty(t *testing.T) {
	capacity := uint64(8)
	rb, err := NewRingBuffer[int](capacity)
	if err != nil {
		t.Fatalf("Failed to create ring buffer: %v", err)
	}

	// Buffer should be empty initially
	if !rb.IsEmpty() {
		t.Fatalf("Ring buffer should be empty initially")
	}

	// Attempt to dequeue from an empty buffer, should fail
	_, err = rb.Dequeue()
	if err == nil {
		t.Fatalf("Expected error when dequeueing from empty buffer, got nil")
	}
	if !errors.Is(err, ErrRingBufferEmpty) {
		t.Fatalf("Expected ErrRingBufferEmpty, got %v", err)
	}

	// Enqueue one item
	if err := rb.Enqueue(42); err != nil {
		t.Fatalf("Failed to enqueue item: %v", err)
	}

	// Buffer should not be empty now
	if rb.IsEmpty() {
		t.Fatalf("Ring buffer should not be empty after enqueueing an item")
	}

	// Dequeue the item
	item, err := rb.Dequeue()
	if err != nil {
		t.Fatalf("Failed to dequeue item: %v", err)
	}
	if item != 42 {
		t.Fatalf("Expected item 42, got %d", item)
	}

	// Buffer should be empty again
	if !rb.IsEmpty() {
		t.Fatalf("Ring buffer should be empty after dequeueing all items")
	}
}

func TestRingBufferConcurrentFull(t *testing.T) {
	capacity := uint64(1024)
	rb, err := NewRingBuffer[int](capacity)
	if err != nil {
		t.Fatalf("Failed to create ring buffer: %v", err)
	}

	var wg sync.WaitGroup
	numProducers := 4
	itemsPerProducer := int(capacity) / numProducers

	// Producers
	for p := 0; p < numProducers; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				for {
					if err := rb.Enqueue(p*itemsPerProducer + i); err == nil {
						break
					} else if errors.Is(err, ErrRingBufferFull) {
						// Buffer is full, wait and retry
						continue
					} else {
						t.Errorf("Unexpected error during enqueue: %v", err)
					}
				}
			}
		}(p)
	}

	// Wait for producers to finish
	wg.Wait()

	// Buffer should be full now
	if !rb.IsFull() {
		t.Fatalf("Ring buffer should be full after producers finish")
	}

	// Attempt to enqueue one more item, should fail
	err = rb.Enqueue(9999)
	if err == nil || !errors.Is(err, ErrRingBufferFull) {
		t.Fatalf("Expected ErrRingBufferFull when enqueueing to full buffer, got %v", err)
	}

	// Dequeue all items
	totalItems := numProducers * itemsPerProducer
	for i := 0; i < totalItems; i++ {
		_, err := rb.Dequeue()
		if err != nil {
			t.Fatalf("Failed to dequeue item %d: %v", i, err)
		}
	}

	// Buffer should be empty now
	if !rb.IsEmpty() {
		t.Fatalf("Ring buffer should be empty after dequeueing all items")
	}
}
