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


func work[V any](f func() V, recv chan <- V, ready chan <- struct{}) {
	defer close(ready)
	defer close(recv)
	recv <- f()
}


func WrapErrorFunc[V any](f func() (V, error)) (func() WithError[V]) {
	return func() WithError[V] {
		v, err := f()
		return WithError[V]{ Value: v, Error: err }
	}
}


func Run[V any](f func() V) *Job[V] {
	recv := make(chan V, 1)
	ready := make(chan struct{})
	consumed := make(chan struct{})

	go work(f, recv, ready)

	return &Job[V]{ recv: recv, ready: ready, consumed: consumed }
}


func (p *Job[V]) Wait() (V, error) {
	return p.WaitContext(context.Background())
}

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
		return v, ctx.Err()
	}
}


func First[V any](jobs ...*Job[V]) (V, error) {
	return FirstContext(context.Background(), jobs...)
}

func FirstContext[V any](ctx context.Context, jobs ...*Job[V]) (V, error) {
	c := make(chan *Job[V])

	cancel := make(chan struct{})
	defer close(cancel)

	n := 0
	for _, job := range jobs {
		select {
		case <- ctx.Done():
			var v V
			return v, ctx.Err()
		case <- job.consumed:
			// Skip already consumed Job.
		default:
			n += 1
			go func(ijob *Job[V]){
				select {
				case <- ijob.ready:
				case <- cancel:
					return
				}

				select {
				case c <- ijob:
				case <- cancel:
				}
			}(job)
		}
	}

	for i := 0; i < n; i++ {
		select {
		case <- ctx.Done():
			var v V
			return v, ctx.Err()
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

	var v V
	return v, ErrAlreadyConsumed
}

func MaybeAll[V any](jobs ...*Job[V]) []WithError[V] {
	return MaybeAllContext(context.Background(), jobs...)
}

func MaybeAllContext[V any](ctx context.Context, jobs ...*Job[V]) []WithError[V] {
	vs := make([]WithError[V], 0, len(jobs))

	for _, job := range jobs {
		v, err := job.WaitContext(ctx)
		vs = append(vs, WithError[V]{ Value: v, Error: err })
	}

	return vs
}
