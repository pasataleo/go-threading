package threading

import (
	"context"
	"reflect"
	"sync"
)

// ImmediateFuture creates a future that runs the given function immediately.
func ImmediateFuture(ctx context.Context, fn func(ctx context.Context)) Future {
	future := &future{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(ctx, reflect.ValueOf(fn), func(values []reflect.Value) {})
	return future
}

// ImmediateFutureV creates a future that runs the given function immediately.
func ImmediateFutureV[V any](ctx context.Context, fn func(ctx context.Context) V) FutureV[V] {
	future := &futureV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(ctx, reflect.ValueOf(fn), func(values []reflect.Value) {
		future.value = values[0].Interface().(V)
	})
	return future
}

// ImmediateFutureE creates a future that runs the given function immediately.
func ImmediateFutureE(ctx context.Context, fn func(ctx context.Context) error) FutureE {
	future := &futureE{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(ctx, reflect.ValueOf(fn), func(values []reflect.Value) {
		if values[0].Interface() != nil {
			future.err = values[0].Interface().(error)
		}
	})
	return future
}

// ImmediateFutureEV creates a future that runs the given function immediately.
func ImmediateFutureEV[V any](ctx context.Context, fn func(ctx context.Context) (V, error)) FutureEV[V] {
	future := &futureEV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	future.immediate(ctx, reflect.ValueOf(fn), func(values []reflect.Value) {
		if values[1].Interface() != nil {
			future.err = values[1].Interface().(error)
		}
		future.value = values[0].Interface().(V)
	})

	return future
}
