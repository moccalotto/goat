package main

import "runtime"

var (
	mainthreadChannel  chan func()
	maxConcurrentCalls = 10
)

func assertInitialized() {
	if mainthreadChannel == nil {
		panic("You mus run your code inside RunInMainthread()")
	}
}

// auto-called when file is loaded
func init() {
	runtime.LockOSThread()
}

func RunInMainthread(fn func()) {
	mainthreadChannel = make(chan func(), maxConcurrentCalls)

	done := make(chan struct{})
	go func() {
		fn()
		done <- struct{}{}
	}()

	for {
		select {
		case f := <-mainthreadChannel:
			f()
		case <-done:
			return
		}
	}
}

// cann a function on the main thread
func Call(fn func()) {
	assertInitialized()
	done := make(chan bool)
	mainthreadChannel <- func() {
		fn()
		done <- true
	}
	<-done
}

// call a function on the main thread, and return any errors it might have returned
func CallErr(fn func() error) error {
	assertInitialized()
	err_chanel := make(chan error)
	mainthreadChannel <- func() {
		err_chanel <- fn()
	}

	return <-err_chanel
}
