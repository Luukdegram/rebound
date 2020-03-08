package thread

import (
	"context"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

var (
	// channel that runs on the main thread. Every call will enter this 'queue'.
	queue = make(chan func())
)

// Run initializes the LockOSThread function and lives until closure of the application.
func Run(context context.Context) {
	var done = false
	for !done {
		select {
		case f := <-queue:
			f()
		case <-context.Done():
			done = true
		}
	}
}

// Call will ensure the the function provided will be run on the mainthread
func Call(f func()) {
	done := make(chan bool, 1)
	queue <- func() {
		f()
		done <- true
	}
	<-done
}

// CallVal will alow you to run a function on the mainthread and return a value
func CallVal(f func() interface{}) interface{} {
	val := make(chan interface{}, 1)
	queue <- func() {
		val <- f()
	}
	return <-val
}

// CallErr is a helper function to call a function on the mainthread that could result in an error
func CallErr(f func() error) error {
	err := make(chan error, 1)
	queue <- func() {
		err <- f()
	}

	return <-err
}

// Locker tells if an object is to be run on the mainthread
type Locker interface {
	IsLocked() bool
}
