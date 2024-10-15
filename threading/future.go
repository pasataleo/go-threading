package threading

import (
	"reflect"
	"sync"
)

// futureImpl abstracts out some common functions for the immediate futures defined in
// this file.
type futureImpl struct {
	// wg is the wait group that is used to wait for the future to finish.
	wg *sync.WaitGroup

	// finished is true if the future is finished, false otherwise.
	finished bool
}

// immediate starts the future with the given function and callback immediately.
//
// The callback is called with the results of the function when the function is finished.
//
// This function should be called by whatever creates the future. Only of immediate or pool should be called.
func (f *futureImpl) immediate(fn reflect.Value, cb func([]reflect.Value)) {
	f.wg.Add(1) // this ensures that the future is not finished until the function is done.

	go func() {
		results := fn.Call(nil)
		cb(results)

		f.finished = true
		f.wg.Done()
	}()
}

// Finished returns true if the future is finished, false otherwise.
func (f *futureImpl) Finished() bool {
	return f.finished
}

// Future is a generic interface that represents a future value.
type Future interface {
	// Get returns the value of the future. If the future is not finished, this method will block until the future is
	// finished.
	Get()

	// Finished returns true if the future is finished, false otherwise.
	Finished() bool
}

var _ Future = (*future)(nil)

type future struct {
	*futureImpl
}

func (f *future) Get() {
	// the basic future just waits for the function to finish
	f.wg.Wait()
}

// FutureV is a generic interface that represents a future value.
type FutureV[V any] interface {
	// Get returns the value of the future. If the future is not finished, this method will block until the future is
	// finished.
	Get() V

	// Finished returns true if the future is finished, false otherwise.
	Finished() bool
}

var _ FutureV[any] = (*futureV[any])(nil)

type futureV[V any] struct {
	*futureImpl

	value V
}

func (f *futureV[V]) Get() V {
	f.wg.Wait()
	return f.value
}

// FutureE is a generic interface that represents a future value.
type FutureE interface {
	// Get returns the value of the future. If the future is not finished, this method will block until the future is
	// finished.
	Get() error

	// Finished returns true if the future is finished, false otherwise.
	Finished() bool
}

var _ FutureE = (*futureE)(nil)

type futureE struct {
	*futureImpl

	err error
}

func (f *futureE) Get() error {
	f.wg.Wait()
	return f.err
}

// FutureEV is a generic interface that represents a future value.
type FutureEV[V any] interface {
	// Get returns the value of the future. If the future is not finished, this method will block until the future is
	// finished.
	Get() (V, error)

	// Finished returns true if the future is finished, false otherwise.
	Finished() bool
}

var _ FutureEV[any] = (*futureEV[any])(nil)

type futureEV[V any] struct {
	*futureImpl

	value V
	err   error
}

func (f *futureEV[V]) Get() (V, error) {
	f.wg.Wait()
	return f.value, f.err
}
