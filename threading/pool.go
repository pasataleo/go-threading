package threading

import (
	"reflect"
	"sync"

	"github.com/pasataleo/go-collections/collections"
	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-objects/objects"
)

var _ objects.Object = (*function)(nil)

type function struct {
	objects.AbstractObject
	f func()
}

// ThreadPool is a thread pool that can be used to run functions concurrently.
type ThreadPool struct {
	// mutex is used to protect the enabled, pending and active pools.
	mutex *sync.Mutex

	// wait is used to wait for all functions to finish.
	wait *sync.WaitGroup

	// enabled is true if the pool is enabled.
	enabled bool

	// size is the number of threads in the pool, this should be a constant.
	size int

	// active counts the number of currently active threads.
	active int

	// pending is a queue of functions that are waiting to immediate.
	pending collections.Queue[*function]

	// ch is the channel that functions are sent to. They are then picked up by the threads executing the work
	// function.
	ch chan func()
}

// NewThreadPool creates a new thread pool with the given number of threads.
func NewThreadPool(size int) *ThreadPool {
	pool := &ThreadPool{
		mutex:   new(sync.Mutex),
		enabled: true,
		active:  0,
		size:    size,
		pending: collections.NewQueue[*function](),
		wait:    new(sync.WaitGroup),
		ch:      make(chan func(), size),
	}

	if size <= 0 {
		panic("size must be greater than 0")
	}

	for i := 0; i < size; i++ {
		// immediate a new goroutine for each thread where they're all just waiting
		// for functions to be enqueued.
		go func() {
			for f := range pool.ch {
				f()
			}
		}()
	}

	return pool
}

// Close closes the thread pool. This will block until all pending functions have been executed.
func (pool *ThreadPool) Close() {
	pool.mutex.Lock()
	pool.enabled = false // mark the pool as closed.
	pool.mutex.Unlock()

	pool.wait.Wait() // now, wait for all the functions to finish.
	close(pool.ch)   // and close the channel.
}

// Running returns true if there are any functions currently running.
func (pool *ThreadPool) Running() bool {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	return pool.active > 0 && !pool.pending.IsEmpty()
}

// Wait waits for all pending functions to finish.
func (pool *ThreadPool) Wait() {
	pool.wait.Wait()
}

// enqueue moves functions from pending into the channel.
//
// enqueue should be called with the pool mutex already locked.
func (pool *ThreadPool) enqueue() {
	for pool.active < pool.size && !pool.pending.IsEmpty() {
		f, _ := pool.pending.Pop()
		pool.active++  // update the number of active threads.
		pool.ch <- f.f // send the function to the channel.
	}
}

// run schedules the given function to run in the thread pool.
func (pool *ThreadPool) run(future *futureImpl, fn reflect.Value, cb func([]reflect.Value)) {
	future.wg.Add(1)

	f := &function{
		f: func() {
			results := fn.Call(nil)
			cb(results)

			// mark the actual future as finished and decrement the wait group.
			future.finished = true
			future.wg.Done()

			// now update the thread pool so it can execute the next one

			pool.mutex.Lock()
			defer pool.mutex.Unlock()

			pool.wait.Done()
			pool.active--
			pool.enqueue()
		},
	}

	pool.pending.Offer(f)
	pool.wait.Add(1)
	pool.enqueue()
}

// Run runs the given function in the given thread pool.
func Run(pool *ThreadPool, f func()) (Future, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if !pool.enabled {
		// if the pool has been closed, return an error.
		return nil, errors.New(nil, ErrorNotActive, "thread pool is not active")
	}

	future := &future{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	pool.run(future.futureImpl, reflect.ValueOf(f), func([]reflect.Value) {})

	return future, nil
}

func RunE(pool *ThreadPool, f func() error) (FutureE, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if !pool.enabled {
		return nil, errors.New(nil, ErrorNotActive, "thread pool is not active")
	}

	future := &futureE{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	pool.run(future.futureImpl, reflect.ValueOf(f), func(values []reflect.Value) {
		future.err = values[0].Interface().(error)
	})

	return future, nil
}

func RunV[V any](pool *ThreadPool, f func() V) (FutureV[V], error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if !pool.enabled {
		return nil, errors.New(nil, ErrorNotActive, "thread pool is not active")
	}

	future := &futureV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	pool.run(future.futureImpl, reflect.ValueOf(f), func(values []reflect.Value) {
		future.value = values[0].Interface().(V)
	})

	return future, nil
}

func RunEV[V any](pool *ThreadPool, f func() (V, error)) (FutureEV[V], error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if !pool.enabled {
		return nil, errors.New(nil, ErrorNotActive, "thread pool is not active")
	}

	future := &futureEV[V]{
		futureImpl: &futureImpl{
			wg: new(sync.WaitGroup),
		},
	}
	pool.run(future.futureImpl, reflect.ValueOf(f), func(values []reflect.Value) {
		future.value = values[0].Interface().(V)
		future.err = values[1].Interface().(error)
	})

	return future, nil
}
