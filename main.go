package main

import (
	"fmt"
	"sync"

	"github.com/keithknott26/easyringbuffer"
)

type App struct {
	sync.Mutex
	Info   *easyringbuffer.StringRingBuffer
	Warn   *easyringbuffer.StringRingBuffer
	Error  *easyringbuffer.StringRingBuffer
	Fatal  *easyringbuffer.StringRingBuffer
	Floats *easyringbuffer.Float64RingBuffer
}

func NewApp() *App {
	return &App{
		//initialize the ring buffers with a length of 100
		Info:   easyringbuffer.NewStringRingBuffer(100),
		Warn:   easyringbuffer.NewStringRingBuffer(100),
		Error:  easyringbuffer.NewStringRingBuffer(100),
		Fatal:  easyringbuffer.NewStringRingBuffer(100),
		Floats: easyringbuffer.NewFloat64RingBuffer(100),
	}
}

func main() {
	// simulate INFO, WARN, ERROR, FATAL log messages
	infoText := "INFO: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	warnText := "WARN: Application has bad references"
	errorText := "ERROR: Line count not found"
	fatalText := "FATAL: Application failed to compile"
	pi := 3.145926

	// initialize the string ring buffer
	app := NewApp()

	// add the log messages to the string ring buffer
	app.Info.Add(infoText)
	app.Warn.Add(warnText)
	app.Error.Add(errorText)
	app.Fatal.Add(fatalText)

	// add the floats onto the float64 ring buffer
	for i := 0; i <= 300; i += 1 {
		app.Floats.Add(float64(i) + pi)
	}

	// print out the last 10 string messages on the buffer
	info := fmt.Sprintf("Last 10 INFO Messages: %s", app.Info.Last(10))
	warn := fmt.Sprintf("Last 10 WARN Messages: %s", app.Warn.Last(10))
	error := fmt.Sprintf("Last 10 ERROR Messages: %s", app.Error.Last(10))
	fatal := fmt.Sprintf("Last 10 FATAL Messages: %s", app.Fatal.Last(10))
	fmt.Println(info)
	fmt.Println(warn)
	fmt.Println(error)
	fmt.Println(fatal)

	// print out the last 10 float64 values on the buffer
	float := fmt.Sprintf("Last 10 Floats: %v", app.Floats.Last(10))
	fmt.Println(float)

	// print the float64 buffer capacity
	fmt.Println(app.Floats.BufCapacity())

	// print all the buffer values
	fmt.Println(app.Floats.BufValues())

	// print all the buffer values between position 1 and 15
	fmt.Println(app.Floats.Slice(1, 15))

	// print all the buffer values between position 10 and 25
	fmt.Println(app.Floats.Position(10, 25))

	// clear all the buffer values
	app.Floats.Reset()

}
