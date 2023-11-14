// Package async implements future / promise pattern.
package async

import (
	"context"
	"errors"
)


type (
	WithError[V any] struct {
		Value V
		Error error
	}

	// Job is a promise object
	// which will receive a result of the asynchronous function.
	// The result can be consumed only once.
	Job[V any] struct {
		recv <- chan V
		ready <- chan struct{}
		consumed chan struct{}
	}
)

var (
	ErrAlreadyConsumed = errors.New("Job has already been consumed")
	ErrReceiverClosed = errors.New("Job receiver channnel has been closed")
)

// WrapErrorFunc returns new wrapped function which returns WithError[V].
func WrapErrorFunc[V any](f func() (V, error)) (func() WithError[V]) {
	return func() WithError[V] {
		v, err := f()
		return WithError[V]{ Value: v, Error: err }
	}
}

func newJob[V any](f func() V) (*Job[V], func()) {
	recv := make(chan V, 1)
	ready := make(chan struct{})
	consumed := make(chan struct{})

	job := Job[V]{ recv: recv, ready: ready, consumed: consumed }
	work := func() {
		defer close(ready)
		defer close(recv)
		recv <- f()
	}
	return &job, work
}

// Run[V] executes function f asynchronously and returns a pointer to the new Job[V].
// The function f must have a single return value of V.
// The function f runs in the new goroutine without any limitations.
// If you want to limit the number of execution simultanously,
// use RunAtWorker[V] function instead.
func Run[V any](f func() V) *Job[V] {
	job, work := newJob(f)

	go work()

	return job
}

// Ready returns a channel which gets signal when the result is ready.
// The result might have been consumed already.
func (p *Job[V]) Ready() <- chan struct{} {
	return p.ready
}

// Consumed returns a channel which gets signal when the result is consumed.
func (p *Job[V]) Consumed() <- chan struct{} {
	return p.consumed
}

// Wait waits the result infinitly. If the result has already been consumed,
// ErrAlreadyConsumed error is returned.
func (p *Job[V]) Wait() (V, error) {
	return p.WaitContext(context.Background())
}

// WaitContext waits the result with Context ctx.
// The ctx doesn't affect running function.
// If the result has already been consumed, ErrAlreadyConsumed error is returned.
// If ctx is cancelled, context.Cause(ctx) error is returned.
func (p *Job[V]) WaitContext(ctx context.Context) (V, error) {
	select {
	case v, ok := <- p.recv:
		if ok {
			close(p.consumed)
			return v, nil
		}
		return v, ErrAlreadyConsumed
	case <- ctx.Done():
		var v V
		return v, context.Cause(ctx)
	}
}

// First waits the first result of jobs infinitly.
// If all jobs have already been consumed, ErrAlreadyConsumed error is returned.
func First[V any](jobs ...*Job[V]) (V, error) {
	return FirstContext(context.Background(), jobs...)
}

// FirstContext waits the first result of jobs with Context ctx.
// The ctx doesn't affect running functions.
// If all jobs have already been consumed, ErrAlreadyConsumed error is returned.
// If ctx is cancelled, context.Cause(ctx) error is returned.
func FirstContext[V any](ctx context.Context, jobs ...*Job[V]) (V, error) {
	c := make(chan *Job[V])

	ctx, cause := context.WithCancelCause(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, job := range jobs {
		select {
		case <- ctx.Done():
			var v V
			return v, context.Cause(ctx)
		case <- job.Consumed():
			// Skip already consumed Job.
		default:
			go func(ijob *Job[V]){
				select {
				case <- ijob.Ready():
				case <- ctx.Done():
					return
				}

				select {
				case c <- ijob:
				case <- ctx.Done():
				}
			}(job)
		}
	}

	go func(){
		for _, job := range jobs {
			select {
			case <- job.Consumed():
			case <- ctx.Done():
				return
			}
		}
		cause(ErrAlreadyConsumed)
	}()

	for {
		select {
		case <- ctx.Done():
			var v V
			return v, context.Cause(ctx)
		case job, ok := <- c:
			if !ok {
				var v V
				return v, ErrReceiverClosed
			}

			v, err := job.WaitContext(ctx)
			if err == nil {
				return v, nil
			}
		}
	}
}

// MaybeAll waits the all results infinitly
// and returns the results as a slice of WithError[V].
// If a job has been already consumed,
// ErrAlreadyConsumed error is set to Error member of WithError[V], otherwise nil.
func MaybeAll[V any](jobs ...*Job[V]) []WithError[V] {
	return MaybeAllContext(context.Background(), jobs...)
}

// MaybeAllContext waits the all results with Context ctx
// and returns the results as a slice of WithError[V].
// The ctx doesn't affect running functions.
// If a job has been already consumed,
// ErrAlreadyConsumed error is set to Error member of WithError[V],
// and if a ctx is cancelled, context.Cause(ctx) error, otherwise nil.
func MaybeAllContext[V any](ctx context.Context, jobs ...*Job[V]) []WithError[V] {
	vs := make([]WithError[V], 0, len(jobs))

	for _, job := range jobs {
		v, err := job.WaitContext(ctx)
		vs = append(vs, WithError[V]{ Value: v, Error: err })
	}

	return vs
}
