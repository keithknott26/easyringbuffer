package easyringbuffer_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/keithknott26/easyringbuffer"
)

func TestNewRingBuffer(t *testing.T) {
	// Test creating a ring buffer with valid capacity
	rb := easyringbuffer.New[int](10)
	if rb == nil {
		t.Fatal("Expected a new ring buffer instance, got nil")
	}
	if rb.Capacity() != 10 {
		t.Fatalf("Expected capacity 10, got %d", rb.Capacity())
	}

	// Test creating a ring buffer with invalid capacity
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when creating ring buffer with capacity 0, but did not panic")
		}
	}()
	easyringbuffer.New
}

func TestEnqueueDequeue(t *testing.T) {
	rb := easyringbuffer.New[int](5)

	// Enqueue items
	for i := 1; i <= 5; i++ {
		err := rb.Enqueue(i)
		if err != nil {
			t.Errorf("Unexpected error on Enqueue: %v", err)
		}
	}

	// Buffer should be full now
	if !rb.IsFull() {
		t.Errorf("Expected buffer to be full")
	}

	// Enqueue should return error now
	err := rb.Enqueue(6)
	if !errors.Is(err, easyringbuffer.ErrRingBufferFull) {
		t.Errorf("Expected ErrRingBufferFull, got %v", err)
	}

	// Dequeue items
	for i := 1; i <= 5; i++ {
		item, err := rb.Dequeue()
		if err != nil {
			t.Errorf("Unexpected error on Dequeue: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}

	// Buffer should be empty now
	if !rb.IsEmpty() {
		t.Errorf("Expected buffer to be empty")
	}

	// Dequeue should return error now
	_, err = rb.Dequeue()
	if !errors.Is(err, easyringbuffer.ErrRingBufferEmpty) {
		t.Errorf("Expected ErrRingBufferEmpty, got %v", err)
	}
}

func TestGetAllAndGetLastN(t *testing.T) {
	rb := easyringbuffer.New[int](5)

	// Enqueue items
	for i := 1; i <= 5; i++ {
		_ = rb.Enqueue(i)
	}

	// Test GetAll
	allItems := rb.GetAll()
	expectedItems := []int{1, 2, 3, 4, 5}
	if len(allItems) != len(expectedItems) {
		t.Errorf("Expected GetAll to return %d items, got %d", len(expectedItems), len(allItems))
	}
	for i, v := range allItems {
		if v != expectedItems[i] {
			t.Errorf("Expected item %d, got %d", expectedItems[i], v)
		}
	}

	// Test GetLastN with n less than size
	last3 := rb.GetLastN(3)
	expectedLast3 := []int{3, 4, 5}
	if len(last3) != len(expectedLast3) {
		t.Errorf("Expected GetLastN to return %d items, got %d", len(expectedLast3), len(last3))
	}
	for i, v := range last3 {
		if v != expectedLast3[i] {
			t.Errorf("Expected item %d, got %d", expectedLast3[i], v)
		}
	}

	// Test GetLastN with n greater than size
	last10 := rb.GetLastN(10)
	if len(last10) != 5 {
		t.Errorf("Expected GetLastN to return %d items, got %d", 5, len(last10))
	}
}

func TestPeek(t *testing.T) {
	rb := easyringbuffer.New[int](5)

	// Peek on empty buffer
	_, err := rb.Peek()
	if !errors.Is(err, easyringbuffer.ErrRingBufferEmpty) {
		t.Errorf("Expected ErrRingBufferEmpty, got %v", err)
	}

	// Enqueue items
	for i := 1; i <= 3; i++ {
		_ = rb.Enqueue(i)
	}

	// Peek should return the first item
	item, err := rb.Peek()
	if err != nil {
		t.Errorf("Unexpected error on Peek: %v", err)
	}
	if item != 1 {
		t.Errorf("Expected Peek to return 1, got %d", item)
	}

	// Dequeue and Peek again
	_, _ = rb.Dequeue()
	item, err = rb.Peek()
	if err != nil {
		t.Errorf("Unexpected error on Peek: %v", err)
	}
	if item != 2 {
		t.Errorf("Expected Peek to return 2, got %d", item)
	}
}

func TestLenAndCapacity(t *testing.T) {
	rb := easyringbuffer.New[int](5)

	if rb.Len() != 0 {
		t.Errorf("Expected Len() == 0, got %d", rb.Len())
	}

	if rb.Capacity() != 5 {
		t.Errorf("Expected Capacity() == 5, got %d", rb.Capacity())
	}

	_ = rb.Enqueue(1)
	_ = rb.Enqueue(2)

	if rb.Len() != 2 {
		t.Errorf("Expected Len() == 2, got %d", rb.Len())
	}

	_, _ = rb.Dequeue()

	if rb.Len() != 1 {
		t.Errorf("Expected Len() == 1, got %d", rb.Len())
	}
}

func TestIsEmptyAndIsFull(t *testing.T) {
	rb := easyringbuffer.New[int](3)

	if !rb.IsEmpty() {
		t.Errorf("Expected IsEmpty() == true")
	}

	if rb.IsFull() {
		t.Errorf("Expected IsFull() == false")
	}

	_ = rb.Enqueue(1)
	if rb.IsEmpty() {
		t.Errorf("Expected IsEmpty() == false")
	}
	if rb.IsFull() {
		t.Errorf("Expected IsFull() == false")
	}

	_ = rb.Enqueue(2)
	if !rb.IsFull() {
		t.Errorf("Expected IsFull() == true")
	}

	_, _ = rb.Dequeue()
	if rb.IsFull() {
		t.Errorf("Expected IsFull() == false")
	}
}

func TestReset(t *testing.T) {
	rb := easyringbuffer.New[int](10)

	for i := 1; i <= 3; i++ {
		_ = rb.Enqueue(i)
	}

	rb.Reset()

	if rb.Len() != 0 {
		t.Errorf("Expected Len() == 0 after Reset(), got %d", rb.Len())
	}

	if !rb.IsEmpty() {
		t.Errorf("Expected IsEmpty() == true after Reset()")
	}

	// Enqueue again after Reset
	_ = rb.Enqueue(10)
	item, err := rb.Dequeue()
	if err != nil {
		t.Errorf("Unexpected error after Reset(): %v", err)
	}
	if item != 10 {
		t.Errorf("Expected item 10 after Reset(), got %d", item)
	}
}

func TestWrapAround(t *testing.T) {
	rb := easyringbuffer.New[int](3)

	_ = rb.Enqueue(1)
	_ = rb.Enqueue(2)
	_ = rb.Enqueue(3)

	_, _ = rb.Dequeue()
	_, _ = rb.Dequeue()

	_ = rb.Enqueue(4)
	_ = rb.Enqueue(5)

	expectedItems := []int{3, 4, 5}
	allItems := rb.GetAll()
	if len(allItems) != len(expectedItems) {
		t.Errorf("Expected %d items, got %d", len(expectedItems), len(allItems))
	}
	for i, v := range allItems {
		if v != expectedItems[i] {
			t.Errorf("Expected item %d, got %d", expectedItems[i], v)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	rb := easyringbuffer.New[int](5000)

	var wg sync.WaitGroup
	numProducers := 5
	numConsumers := 5
	itemsPerProducer := 1000

	// Producers
	for p := 0; p < numProducers; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				for {
					err := rb.Enqueue(p*itemsPerProducer + i)
					if err == nil {
						break
					}
					if !errors.Is(err, easyringbuffer.ErrRingBufferFull) {
						t.Errorf("Unexpected error on Enqueue: %v", err)
						return
					}
					// Buffer is full, retry
				}
			}
		}(p)
	}

	// Consumers
	totalItems := numProducers * itemsPerProducer
	consumedCount := 0
	consumedItems := make([]int, 0, totalItems)
	mu := sync.Mutex{}

	for c := 0; c < numConsumers; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				item, err := rb.Dequeue()
				if err == nil {
					mu.Lock()
					consumedItems = append(consumedItems, item)
					consumedCount++
					mu.Unlock()
				} else if errors.Is(err, easyringbuffer.ErrRingBufferEmpty) {
					// Buffer is empty, check if producers are done
					mu.Lock()
					done := consumedCount >= totalItems
					mu.Unlock()
					if done {
						return
					}
					// Buffer is empty, retry
				} else {
					t.Errorf("Unexpected error on Dequeue: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Verify all items were consumed
	if consumedCount != totalItems {
		t.Errorf("Expected %d items consumed, got %d", totalItems, consumedCount)
	}

	// Optionally, check for duplicates or missing items
	itemMap := make(map[int]bool)
	for _, item := range consumedItems {
		if itemMap[item] {
			t.Errorf("Duplicate item found: %d", item)
		}
		itemMap[item] = true
	}

	if len(itemMap) != totalItems {
		t.Errorf("Expected %d unique items, got %d", totalItems, len(itemMap))
	}
}

func TestRingBufferGeneric(t *testing.T) {
	// Test with strings
	rb := easyringbuffer.New[string](3)
	_ = rb.Enqueue("one")
	_ = rb.Enqueue("two")
	_ = rb.Enqueue("three")

	item, err := rb.Dequeue()
	if err != nil {
		t.Errorf("Unexpected error on Dequeue: %v", err)
	}
	if item != "one" {
		t.Errorf("Expected 'one', got '%s'", item)
	}

	// Test with custom struct
	type MyStruct struct {
		ID   int
		Name string
	}
	rbStruct := easyringbuffer.New[MyStruct](3)
	_ = rbStruct.Enqueue(MyStruct{ID: 1, Name: "Alice"})
	_ = rbStruct.Enqueue(MyStruct{ID: 2, Name: "Bob"})

	structItem, err := rbStruct.Dequeue()
	if err != nil {
		t.Errorf("Unexpected error on Dequeue: %v", err)
	}
	if structItem.ID != 1 || structItem.Name != "Alice" {
		t.Errorf("Expected ID=1, Name='Alice', got ID=%d, Name='%s'", structItem.ID, structItem.Name)
	}
}

func TestEnqueueNil(t *testing.T) {
	// Note: Only applicable if T can be a pointer type
	rb := easyringbuffer.New[*int](1)

	var ptr *int = nil
	err := rb.Enqueue(ptr)
	if err != nil {
		t.Errorf("Unexpected error on Enqueue: %v", err)
	}

	item, err := rb.Dequeue()
	if err != nil {
		t.Errorf("Unexpected error on Dequeue: %v", err)
	}
	if item != nil {
		t.Errorf("Expected nil item, got %v", item)
	}
}
