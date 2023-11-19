// queuechan package
//
// Queue based infinite length channel
package queuechan

import (
	"context"
)

// Create new channel set
//
// # Arguments
// * `ctx`: `context.Context` - Context to stop Queue Channnel
//
// # Returns
// * `in`: `chan <- T` - Input channel to Queue
// * `out`: `<- chan T` - Output channel from Queue
//
// # Context
// * `Done()` -> Don't take out from `in` any more, but do put to `in` for cleanup.
func New[T any](ctx context.Context) (chan <- T, <- chan T) {
	in := make(chan T, 0)
	out := make(chan T, 0)

	go func(in <- chan T, out chan <- T){
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

		// Clean up
		for len(queue) > 0 {
			out <- queue[0]
			queue = queue[1:]
		}
	}(in, out)

	return in, out
}
