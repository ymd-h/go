// Package qchan provides queue based infinit length channel
package qchan

import (
	"context"
	"errors"
)

type (
	// Queue[T] is a queue based infinif length channel.
	Queue[T any] struct {
		in chan <- T
		out <- chan T
		ctx context.Context
	}
)


var (
	ErrInputClosed = errors.New("Input channel has aleady been closed.")
	ErrUnknown = errors.New("Finish with unknown reason")
)

// New[T] creates a new Queue[T] and returns a point to it.
func New[T any]() *Queue[T] {
	return NewWithContext[T](context.Background())
}

// NewWithContext[T] creates a new Queue[T] and returns a pointer to it.
// If ctx is cancelled, Queue[T] will not consumed input channel,
// and the input channel will be blocked, however,
// remained values still will be put into output channel.
func NewWithContext[T any](ctx context.Context) *Queue[T] {
	in := make(chan T, 0)
	out := make(chan T, 0)
	ctx, cause := context.WithCancelCause(ctx)

	go func(in <- chan T, out chan <- T){
		defer cause(ErrUnknown)
		defer close(out)
		queue := make([]T, 0)

		LOOP:
		for {
			if len(queue) == 0 {
				select {
				case <- ctx.Done():
					return
				case v, ok := <- in:
					if !ok {
						// `in` is closed
						// No cleanup is needed
						cause(ErrInputClosed)
						return
					}
					queue = append(queue, v)
				}
			}

			select {
			case <- ctx.Done():
				break LOOP
			case v, ok := <- in:
				if ok {
					queue = append(queue, v)
				} else {
					// `in` is closed
					// Go to cleanup
					cause(ErrInputClosed)
					break LOOP
				}
			case out <- queue[0]:
				queue = queue[1:]
			}
		}

		// ctx shoule have already been cancelled, however, we ensure canncel it.
		cause(ErrUnknown)

		// Clean up
		for len(queue) > 0 {
			out <- queue[0]
			queue = queue[1:]
		}
	}(in, out)

	return &Queue[T]{ in: in, out: out, ctx: ctx }
}

// In returns input channel.
func (q *Queue[T]) In() chan <- T {
	return q.in
}

// Out returns output channel.
func (q *Queue[T]) Out() <- chan T {
	return q.out
}

// Done returns done channel, which will be closed
// when Queue[T] stops consuming its input channel.
func (q *Queue[T]) Done() <- chan struct{} {
	return q.ctx.Done()
}

// Error returns error explaining cancell reason.
func (q *Queue[T]) Error() error {
	return context.Cause(q.ctx)
}
