package threading

import (
	"sync"
)

// ThreadPool is a pool of workers that can be used to run functions concurrently.
//
// A thread pool will block the current thread until a worker is available to run the function.
type ThreadPool struct {
	workerLock, waitLock sync.Mutex
	waiter               sync.WaitGroup
	current, workers     int
}

func NewThreadPool(workers int) *ThreadPool {
	tp := &ThreadPool{
		workers: workers,
	}
	return tp
}

// waitForWorker will block the current thread until a worker is available.
func (tp *ThreadPool) reserveWorker() {
	tp.waitLock.Lock() // Only one thread can wait for a worker at a time.
	defer tp.waitLock.Unlock()

	tp.workerLock.Lock()
	if tp.current >= tp.workers {
		tp.workerLock.Unlock()
		tp.waiter.Wait()
	} else {
		tp.workerLock.Unlock()
	}

	tp.workerLock.Lock()
	tp.current++
	if tp.current >= tp.workers {
		tp.waiter.Add(1)
	}
	tp.workerLock.Unlock()
}

func (tp *ThreadPool) releaseWorker() {
	tp.workerLock.Lock()
	if tp.current >= tp.workers {
		tp.waiter.Done()
	}
	tp.current--
	tp.workerLock.Unlock()
}

func Run(tp *ThreadPool, fn func()) Future {
	tp.reserveWorker() // Wait for a worker to be available.

	return ImmediateFuture(func() {
		// compute the value in the future.
		fn()
		tp.releaseWorker()
	})
}

func RunV[V any](tp *ThreadPool, fn func() V) FutureV[V] {
	tp.reserveWorker() // Wait for a worker to be available.

	return ImmediateFutureV[V](func() V {
		// compute the value in the future.
		value := fn()
		tp.releaseWorker()
		return value
	})
}

func RunEV[V any](tp *ThreadPool, fn func() (V, error)) FutureEV[V] {
	tp.reserveWorker() // Wait for a worker to be available.

	return ImmediateFutureEV[V](func() (V, error) {
		// compute the value in the future.
		value, err := fn()
		tp.releaseWorker()
		return value, err
	})
}

func RunE(tp *ThreadPool, fn func() error) FutureE {
	tp.reserveWorker() // Wait for a worker to be available.

	return ImmediateFutureE(func() error {
		// compute the value in the future.
		err := fn()
		tp.releaseWorker()
		return err
	})
}
