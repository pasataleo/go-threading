package threading

import (
	"reflect"
	"sync"
)

// ImmediateFuture creates a future that runs the given function immediately.
func ImmediateFuture(fn func()) Future {
	future := &future{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(reflect.ValueOf(fn), func(values []reflect.Value) {})
	return future
}

// ImmediateFutureV creates a future that runs the given function immediately.
func ImmediateFutureV[V any](fn func() V) FutureV[V] {
	future := &futureV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(reflect.ValueOf(fn), func(values []reflect.Value) {
		future.value = values[0].Interface().(V)
	})
	return future
}

// ImmediateFutureE creates a future that runs the given function immediately.
func ImmediateFutureE(fn func() error) FutureE {
	future := &futureE{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(reflect.ValueOf(fn), func(values []reflect.Value) {
		if values[0].Interface() != nil {
			future.err = values[0].Interface().(error)
		}
	})
	return future
}

// ImmediateFutureEV creates a future that runs the given function immediately.
func ImmediateFutureEV[V any](fn func() (V, error)) FutureEV[V] {
	future := &futureEV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(reflect.ValueOf(fn), func(values []reflect.Value) {
		if values[1].Interface() != nil {
			future.err = values[1].Interface().(error)
		}
		future.value = values[0].Interface().(V)
	})

	return future
}
