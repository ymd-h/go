package async

import (
	"context"
)

type (
	IWorker interface {
		Send(context.Context, func()) error
	}

	Worker struct {
		send chan <- func()
	}
)

func NewWorker(n uint) *Worker {
	c := make(chan func())

	for i := uint(0); i < n; i++ {
		go func(){
			for {
				f, ok := <- c
				if !ok {
					return
				}
				f()
			}
		}()
	}

	return &Worker{
		send: c,
	}
}

func NewLazyWorker(n uint) *Worker {
	c := make(chan func())
	sem := make(chan struct{}, n)

	go func(){
		for {
			f, ok := <- c
			if !ok {
				return
			}

			sem <- struct{}{}
			go func(){
				defer func(){ <-sem }()
				f()
			}()
		}
	}()

	return &Worker{
		send: c,
	}
}

func (w *Worker) Send(ctx context.Context, f func()) error {
	select {
	case w.send <- f:
		return nil
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
