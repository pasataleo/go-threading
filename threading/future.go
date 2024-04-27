package threading

type Future interface {
	Finished() bool
	Get()
}

type FutureV[V any] interface {
	Finished() bool
	Get() V
}

type FutureE interface {
	Finished() bool
	Get() error
}

type FutureEV[V any] interface {
	Finished() bool
	Get() (V, error)
}
