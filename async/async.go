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

func Run[V any](f func() V) *Job[V] {
	job, work := newJob(f)

	go work()

	return job
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

	ctx, cause := context.WithCancelCause(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, job := range jobs {
		select {
		case <- ctx.Done():
			var v V
			return v, ctx.Err()
		case <- job.consumed:
			// Skip already consumed Job.
		default:
			go func(ijob *Job[V]){
				select {
				case <- ijob.ready:
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
			case <- job.consumed:
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
