package async

import (
	"context"
	"errors"
)

type (
	// IWorker is a Worker interface
	IWorker interface {
		Send(context.Context, func()) error
	}

	// Worker is a job worker which might limit the number of goroutine
	Worker struct {
		send chan <- func()
		done <- chan struct{}
	}
)

var (
	ErrAlreadyShutdown = errors.New("Worker has already been shut down")
)

// NewWorker creates new worker goroutines and returns a pointer to the new Worker.
// The Context ctx is used to stop the Worker.
// n is the number of goroutine to be prepared.
func NewWorker(ctx context.Context, n uint) *Worker {
	c := make(chan func())

	d := make([]chan struct{}, 0, n)
	for i := uint(0); i < n; i++ {
		d = append(d, make(chan struct{}))
	}

	for _, di := range d {
		go func(dd chan <- struct{}){
			defer close(dd)
			for {
				select {
				case f, ok := <- c:
					if !ok {
						return
					}
					f()
				case <- ctx.Done():
					return
				}
			}
		}(di)
	}

	done := make(chan struct{})
	go func(){
		defer close(done)
		for _, di := range d {
			<- di
		}
	}()

	return &Worker{
		send: c,
		done: done,
	}
}

// NewLazyWorker creates helper goroutine and returns a pointer to the new Worker.
// The helper goroutine will lazily create a worker goroutine
// when the Worker receives a job request.
// The worker goroutine will terminate after finishing the job.
// (It is not necessary to consume the result.)
// The Context ctx is used to stop the Worker.
// The number of worker goroutine is limited by n.
func NewLazyWorker(ctx context.Context, n uint) *Worker {
	c := make(chan func())
	sem := make(chan struct{}, n)
	done := make(chan struct{})

	go func(){
		defer close(done)
		for {
			select {
			case sem <- struct{}{}:
			case <- ctx.Done():
				return
			}

			select {
			case f, ok := <- c:
				if !ok {
					return
				}

				go func(){
					defer func(){ <-sem }()
					f()
				}()
			case <- ctx.Done():
				return
			}
		}
	}()

	return &Worker{
		send: c,
		done: done,
	}
}


// Send sends the prepared job function to the Worker.
// Worker might block and Context ctx can cancel it.
// If the Worker has already been shutdown, ErrAlreadyShutdown is returned,
// and if ctx is canceled, context.Cause(ctx) error, otherwise nil.
//
// This method is not intended to call directly, but to be used in RunAtWorker[V].
// In order to implement custom worker class, Send is public method.
func (w *Worker) Send(ctx context.Context, f func()) error {
	select {
	case w.send <- f:
		return nil
	case <- w.done:
		return ErrAlreadyShutdown
	case <- ctx.Done():
		return context.Cause(ctx)
	}
}


// RunAtWorker[V] executes function f at IWorker w asynchronously,
// and returns a pointer to the new Job[V].
// Context ctx is used to cancel sending the job to the worker,
// it doesn't affect the worker or the job.
// If w.Send returns error, the error is returned.
func RunAtWorker[V any](ctx context.Context, w IWorker, f func() V) (*Job[V], error) {
	job, work := newJob(f)

	if err := w.Send(ctx, work); err != nil {
		return nil, err
	}

	return job, nil
}
