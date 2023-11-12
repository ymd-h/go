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

func NewWorker(n uint) *Worker {
	c := make(chan func())
	done := make(chan struct{})

	d := make([]chan struct{}, n)
	for _, di := range d {
		go func(dd chan <- struct{}){
			defer close(dd)
			for {
				f, ok := <- c
				if !ok {
					return
				}
				f()
			}
		}(di)
	}

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

func NewLazyWorker(n uint) *Worker {
	c := make(chan func())
	sem := make(chan struct{}, n)
	done := make(chan struct{})

	go func(){
		defer close(done)
		for {
			sem <- struct{}{}
			f, ok := <- c
			if !ok {
				return
			}

			go func(){
				defer func(){ <-sem }()
				f()
			}()
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
		return ctx.Err()
	}
}


func RunAtWorker[V any](ctx context.Context, w IWorker, f func() V) (*Job[V], error) {
	c := make(chan (chan <- V))
	done := make(chan struct{})

	if err := w.Send(ctx, func(){ work(f, c, done) }); err != nil {
		return nil, err
	}

	return &Job[V]{send: c, done: done}, nil
}
