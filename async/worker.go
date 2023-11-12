package async

import (
	"context"
	"errors"
)

type (
	IWorker interface {
		Send(context.Context, func()) error
	}

	Worker struct {
		send chan <- func()
		done <- chan struct{}
	}
)

var (
	ErrAlreadyShutdown = errors.New("Worker has already been shut down")
)

func NewWorker(ctx context.Context, n uint) *Worker {
	c := make(chan func())

	d := make([]chan struct{}, n)
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


func RunAtWorker[V any](ctx context.Context, w IWorker, f func() V) (*Job[V], error) {
	job, work := newJob(f)

	if err := w.Send(ctx, work); err != nil {
		return nil, err
	}

	return job, nil
}
