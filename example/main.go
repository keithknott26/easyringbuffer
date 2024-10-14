package main

import (
	"fmt"

	"github.com/keithknott26/easyringbuffer/ringbuffer"
)

type App struct {
	Info   *ringbuffer.RingBuffer[string]
	Warn   *ringbuffer.RingBuffer[string]
	Error  *ringbuffer.RingBuffer[string]
	Fatal  *ringbuffer.RingBuffer[string]
	Floats *ringbuffer.RingBuffer[float64]
}

func NewApp() *App {
	// Initialize ring buffers with a capacity of 128 (must be a power of two)
	infoBuffer, _ := ringbuffer.NewRingBuffer
	warnBuffer, _ := ringbuffer.NewRingBuffer
	errorBuffer, _ := ringbuffer.NewRingBuffer
	fatalBuffer, _ := ringbuffer.NewRingBuffer
	floatBuffer, _ := ringbuffer.NewRingBuffer

	return &App{
		Info:   infoBuffer,
		Warn:   warnBuffer,
		Error:  errorBuffer,
		Fatal:  fatalBuffer,
		Floats: floatBuffer,
	}
}

func main() {
	// Simulate INFO, WARN, ERROR, FATAL log messages
	infoText := "INFO: Lorem ipsum dolor sit amet..."
	warnText := "WARN: Application has bad references"
	errorText := "ERROR: Line count not found"
	fatalText := "FATAL: Application failed to compile"
	pi := 3.1415926

	// Initialize the app with ring buffers
	app := NewApp()

	// Add the log messages to the ring buffers
	app.Info.Enqueue(infoText)
	app.Warn.Enqueue(warnText)
	app.Error.Enqueue(errorText)
	app.Fatal.Enqueue(fatalText)

	// Add float values to the float64 ring buffer
	for i := 0; i <= 300; i++ {
		app.Floats.Enqueue(float64(i) + pi)
	}

	// Print out the last 10 INFO messages
	infoMessages := getLastN(app.Info, 10)
	fmt.Printf("Last %d INFO Messages: %v\n", len(infoMessages), infoMessages)

	// Similarly for other log levels
	warnMessages := getLastN(app.Warn, 10)
	fmt.Printf("Last %d WARN Messages: %v\n", len(warnMessages), warnMessages)

	errorMessages := getLastN(app.Error, 10)
	fmt.Printf("Last %d ERROR Messages: %v\n", len(errorMessages), errorMessages)

	fatalMessages := getLastN(app.Fatal, 10)
	fmt.Printf("Last %d FATAL Messages: %v\n", len(fatalMessages), fatalMessages)

	// Print out the last 10 float64 values
	floatValues := getLastN(app.Floats, 10)
	fmt.Printf("Last %d Floats: %v\n", len(floatValues), floatValues)

	// Print the float64 buffer capacity and size
	fmt.Printf("Float Buffer Capacity: %d\n", app.Floats.Capacity())
	fmt.Printf("Float Buffer Size: %d\n", app.Floats.Size())

	// Print all buffer values (may be large)
	allFloats := getAllValues(app.Floats)
	fmt.Printf("All Floats: %v\n", allFloats)

	// Print buffer values between positions 1 and 15
	slice1 := getSlice(app.Floats, 1, 15)
	fmt.Printf("Floats from position 1 to 15: %v\n", slice1)

	// Print buffer values between positions 10 and 25
	slice2 := getSlice(app.Floats, 10, 25)
	fmt.Printf("Floats from position 10 to 25: %v\n", slice2)

	// Clear all buffer values
	app.Floats.Reset()
	fmt.Println("Float buffer reset.")

	// Ensure the buffer is empty
	if app.Floats.IsEmpty() {
		fmt.Println("Float buffer is now empty.")
	}
}

// Helper functions to retrieve values from the ring buffer

// getLastN retrieves the last n items from the ring buffer
func getLastN[T any](rb *ringbuffer.RingBuffer[T], n int) []T {
	size := int(rb.Size())
	if n > size {
		n = size
	}
	result := make([]T, 0, n)
	for i := size - n; i < size; i++ {
		item, _ := rb.GetAt(i)
		result = append(result, item)
	}
	return result
}

// getAllValues retrieves all items from the ring buffer
func getAllValues[T any](rb *ringbuffer.RingBuffer[T]) []T {
	size := int(rb.Size())
	result := make([]T, 0, size)
	for i := 0; i < size; i++ {
		item, _ := rb.GetAt(i)
		result = append(result, item)
	}
	return result
}

// getSlice retrieves items between start and end positions
func getSlice[T any](rb *ringbuffer.RingBuffer[T], start, end int) []T {
	size := int(rb.Size())
	if start < 0 {
		start = 0
	}
	if end > size {
		end = size
	}
	if start >= end {
		return []T{}
	}
	result := make([]T, 0, end-start)
	for i := start; i < end; i++ {
		item, _ := rb.GetAt(i)
		result = append(result, item)
	}
	return result
}
