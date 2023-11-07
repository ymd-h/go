package async

import (
	"context"
)

type (
	IWorker interface {
		Send(func())
		SendWithContext(context.Context, func()) bool
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

func (w *Worker) Send(f func()) {
	w.SendWithContext(context.Background())
}

func (w *Worker) SendWithContext(ctx context.Context, f func()) bool {
	select {
	case w.send <- f:
		return true
	case <- ctx.Done():
		return false
	}
}


func RunAtWorker[V any](w IWorker, f func() V) *Job[V] {
	c := make(chan (chan <- V))
	done := make(chan struct{})

	w.Send(func(){
		work(f, c, done)
	})

	return &Job[V]{send: c, done: done}
}

func RunWithErrorAtWorker[V any](w IWorker, f func() (V, error)) *Job[WithError[V]] {
	return RunAtWorker(w, wrapWithError(f))
}
