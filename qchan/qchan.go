// Package qchan provides queue based infinit length channel
package qchan

import (
	"context"
)

type (
	// Queue[T] is a queue based infinif length channel.
	Queue[T any] struct {
		in chan <- T
		out <- chan T
		done <- chan struct{}
	}
)

// New[T] creates new Queue[T] and returns a pointer to it.
// If ctx is cancelled, Queue[T] will not consumed input channel,
// and the input channel will be blocked, however,
// remained values still will be put into output channel.
func New[T any](ctx context.Context) *Queue[T] {
	in := make(chan T, 0)
	out := make(chan T, 0)
	done, cancel := context.WithCancel(context.Background())

	go func(in <- chan T, out chan <- T){
		defer cancel()
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
					break LOOP
				}
			case out <- queue[0]:
				queue = queue[1:]
			}
		}

		// Cancel since Queue[T] will not consume input channel any more.
		cancel()

		// Clean up
		for len(queue) > 0 {
			out <- queue[0]
			queue = queue[1:]
		}
	}(in, out)

	return &Queue{ in: in, out: out, done: done.Done() }
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
	return q.done
}
