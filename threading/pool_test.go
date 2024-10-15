package threading

import (
	"sync"
	"testing"
	"time"

	"github.com/pasataleo/go-testing/tests"
)

func TestThreadPool(t *testing.T) {
	pool := NewThreadPool(1)

	data := 0
	tests.Execute2E(Run(pool, func() {
		data = 1
	})).NoError(t)

	pool.Close()

	tests.Execute(data).Equal(t, 1)
}

func TestThreadPool_multi(t *testing.T) {
	pool := NewThreadPool(5)

	mutex := new(sync.Mutex)
	keys := make(map[int]bool)
	fn := func(key int) func() {
		return func() {
			t.Logf("starting %d", key)
			time.Sleep(wait)

			mutex.Lock()
			keys[key] = true
			mutex.Unlock()

			t.Logf("ending %d", key)
		}
	}

	for ix := 0; ix < 10; ix++ {
		tests.Execute2E(Run(pool, fn(ix))).NoError(t)
	}
	pool.Close()

	for ix := 0; ix < 10; ix++ {
		tests.Execute(keys[ix]).Equal(t, true)
	}
}
