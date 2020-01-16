package easyringbuffer

import (
	"sync"
)

type Float64RingBuffer struct {
	sync.Mutex
	Buffer   []float64
	Length   int
	Capacity int
	Tail     int
}

func NewFloat64RingBuffer(capacity int) *Float64RingBuffer {
	return &Float64RingBuffer{
		Buffer:   make([]float64, capacity, capacity),
		Length:   0,
		Capacity: capacity,
		Tail:     0,
	}
}

func (r *Float64RingBuffer) BufValues() []float64 {
	return r.Buffer
}
func (r *Float64RingBuffer) Reset() {
	r.Lock()
	defer r.Unlock()
	r.Tail = 0
	r.Length = 0
}
func (r *Float64RingBuffer) Len() int {
	return r.Length
}

func (r *Float64RingBuffer) BufCapacity() int {
	return r.Capacity
}

func (r *Float64RingBuffer) Add(v float64) {
	r.Lock()
	if r.Length < r.Capacity {
		r.Length += 1
	}
	r.Buffer[r.Tail] = v
	r.Tail = (r.Tail + 1) % r.Capacity
	r.Unlock()
}

func (r *Float64RingBuffer) Slice(i, j int) []float64 {
	r.Lock()
	defer r.Unlock()

	if r.Length < r.Capacity {
		if r.Length < j {
			j = r.Length
		}
		return r.Buffer[i:j]
	}
	s := append(r.Buffer[r.Tail:r.Capacity], r.Buffer[:r.Tail]...)
	return s[i:j]
}

func (r *Float64RingBuffer) Position(i, j int) []float64 {
	r.Lock()
	defer r.Unlock()
	// end exceeds length
	if j > r.Length {
		j = r.Length
	}
	// i exceeds length
	if i > r.Length {
		i = 0
		j = r.Length
	}
	// start is less than 0
	if i < 0 {
		i = 0
	}
	// end is less than 0
	if j < 0 {
		i = 0
		j = r.Length
	}
	if r.Length < r.Capacity {
		if r.Length < j {
			j = r.Length
		}
		return r.Buffer[i:j]
	}
	s := append(r.Buffer[r.Tail:r.Capacity], r.Buffer[:r.Tail]...)
	return s[i:j]
}

func (r *Float64RingBuffer) Last(n int) []float64 {
	r.Lock()
	start := r.Length - n
	if start < 0 {
		start = 0
	}
	r.Unlock()
	return r.Slice(start, r.Length)
}

type StringRingBuffer struct {
	sync.Mutex
	Buffer   []string
	Length   int
	Capacity int
	Tail     int
}

func NewStringRingBuffer(capacity int) *StringRingBuffer {
	return &StringRingBuffer{
		Buffer:   make([]string, capacity, capacity),
		Length:   0,
		Capacity: capacity,
		Tail:     0,
	}
}

func (r *StringRingBuffer) BufValues() []string {
	return r.Buffer
}

func (r *StringRingBuffer) Reset() {
	r.Lock()
	defer r.Unlock()
	r.Tail = 0
	r.Length = 0
}
func (r *StringRingBuffer) Len() int {
	return r.Length
}

func (r *StringRingBuffer) BufCapacity() int {
	return r.Capacity
}

func (r *StringRingBuffer) Add(v string) {
	r.Lock()
	if r.Length < r.Capacity {
		r.Length += 1
	}
	r.Buffer[r.Tail] = v
	r.Tail = (r.Tail + 1) % r.Capacity
	r.Unlock()
}

func (r *StringRingBuffer) Slice(i, j int) []string {
	r.Lock()
	defer r.Unlock()
	if r.Length < r.Capacity {
		if r.Length < j {
			j = r.Length
		}
		return r.Buffer[i:j]
	}
	s := append(r.Buffer[r.Tail:r.Capacity], r.Buffer[:r.Tail]...)
	return s[i:j]
}

func (r *StringRingBuffer) Position(i, j int) []string {
	r.Lock()
	defer r.Unlock()
	// end exceeds length
	if j > r.Length {
		j = r.Length
	}
	// i exceeds length
	if i > r.Length {
		i = 0
		j = r.Length
	}
	// start is less than 0
	if i < 0 {
		i = 0
	}
	// end is less than 0
	if j < 0 {
		i = 0
		j = r.Length
	}

	if r.Length < r.Capacity {
		if r.Length < j {
			j = r.Length
		}
		return r.Buffer[i:j]
	}
	s := append(r.Buffer[r.Tail:r.Capacity], r.Buffer[:r.Tail]...)
	return s[i:j]
}

func (r *StringRingBuffer) Last(n int) []string {
	r.Lock()
	start := r.Length - n
	if start < 0 {
		start = 0
	}
	r.Unlock()
	return r.Slice(start, r.Length)
}
