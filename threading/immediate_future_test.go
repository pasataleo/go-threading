package threading

import (
	"context"
	"testing"
	"time"

	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-testing/tests"
)

var (
	wait = 10 * time.Millisecond
)

func TestImmediateFuture(t *testing.T) {
	data := 0
	f := ImmediateFuture(context.Background(), func(ctx context.Context) {
		time.Sleep(wait)
		data = 1
	})
	f.Get()

	tests.Execute(data).Equal(t, 1)
}

func TestImmediateFutureV(t *testing.T) {
	f := ImmediateFutureV(context.Background(), func(ctx context.Context) int {
		time.Sleep(wait)
		return 1
	})
	tests.Execute(f.Get()).Equal(t, 1)
}

func TestImmediateFutureE(t *testing.T) {
	f := ImmediateFutureE(context.Background(), func(ctx context.Context) error {
		time.Sleep(wait)
		return nil
	})
	tests.ExecuteE(f.Get()).NoError(t)
}

func TestImmediateFutureE_errors(t *testing.T) {
	f := ImmediateFutureE(context.Background(), func(ctx context.Context) error {
		time.Sleep(wait)
		return errors.New(nil, errors.ErrorCodeUnknown, "error")
	})
	tests.ExecuteE(f.Get()).MatchesError(t, "error")
}

func TestImmediateFutureEV(t *testing.T) {
	f := ImmediateFutureEV(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(wait)
		return 1, nil
	})
	tests.Execute2E(f.Get()).NoError(t).Equal(t, 1)
}

func TestImmediateFutureEV_errors(t *testing.T) {
	f := ImmediateFutureEV(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(wait)
		return 0, errors.New(nil, errors.ErrorCodeUnknown, "error")
	})
	tests.Execute2E(f.Get()).MatchesError(t, "error").Equal(t, 0)
}
