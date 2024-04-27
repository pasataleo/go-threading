package threading

import "sync"

//
// Simple future implementation.
//

type immediateFuture struct {
	sync.WaitGroup

	finished bool
}

func ImmediateFuture(fn func()) Future {
	future := &immediateFuture{}
	future.Add(1)

	go func() {
		fn()
		future.finished = true
		future.Done()
	}()

	return future
}

func (f *immediateFuture) Finished() bool {
	return f.finished
}

func (f *immediateFuture) Get() {
	f.Wait()
}

//
// Future with a value.
//

type immediateFutureV[V any] struct {
	sync.WaitGroup

	value    V
	finished bool
}

func ImmediateFutureV[V any](fn func() V) FutureV[V] {
	future := &immediateFutureV[V]{}
	future.Add(1)

	go func() {
		future.value = fn()
		future.finished = true
		future.Done()
	}()

	return future
}

func (f *immediateFutureV[V]) Finished() bool {
	return f.finished
}

func (f *immediateFutureV[V]) Get() V {
	f.Wait()
	return f.value
}

//
// Future with an error.
//

type immediateFutureE struct {
	sync.WaitGroup

	err      error
	finished bool
}

func ImmediateFutureE(fn func() error) FutureE {
	future := &immediateFutureE{}
	future.Add(1)

	go func() {
		future.err = fn()
		future.finished = true
		future.Done()
	}()

	return future
}

func (f *immediateFutureE) Finished() bool {
	return f.finished
}

func (f *immediateFutureE) Get() error {
	f.Wait()
	return f.err
}

//
// Future with a value and an error.
//

type immediateFutureEV[V any] struct {
	sync.WaitGroup

	value    V
	err      error
	finished bool
}

func ImmediateFutureEV[V any](fn func() (V, error)) FutureEV[V] {
	future := &immediateFutureEV[V]{}
	future.Add(1)

	go func() {
		future.value, future.err = fn()
		future.finished = true
		future.Done()
	}()

	return future
}

func (f *immediateFutureEV[V]) Finished() bool {
	return f.finished
}

func (f *immediateFutureEV[V]) Get() (V, error) {
	f.Wait()
	return f.value, f.err
}
