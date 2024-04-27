package threading

// ThreadQueue is a pool of workers that can be used to run functions concurrently.
//
// A thread queue always returns immediately, while the function will pause in the background until a worker is
// available to run the function.
type ThreadQueue struct {
}

func NewThreadQueue(workers int) *ThreadQueue {
	tq := &ThreadQueue{}
	return tq
}
