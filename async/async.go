package async

import (
	"context"
	"errors"
	"fmt"
)


type (
	WithError[V any] struct {
		Value V
		Error error
	}

	Job[V any] struct {
		send chan <- (chan <- V)
		done <- chan struct{}
	}
)

var (
	ErrAlreadyDone = errors.New("Job has already been done")
	ErrSendReceiver = fmt.Errorf("Fail to send receiver channel: %w", ErrAlreadyDone)
	ErrReceiverClosed = errors.New("Job receiver channnel has been closed")
)


func work[V any](f func() V, c <- chan (chan <- V), done chan <- struct{}) {
	defer close(done)
	v := f()

	conn, ok := <- c
	if !ok {
		return
	}

	for {
		select {
		case conn <- v:
			return
		case recv, ok := <- c:
			if ok {
				conn = recv
			} else {
				// 'c' is closed
				conn <- v
				return
			}
		}
	}
}

func WrapErrorFunc[V any](f func() (V, error)) (func() WithError[V]) {
	return func() WithError[V] {
		v, err := f()
		return WithError[V]{ Value: v, Error: err }
	}
}


func Run[V any](f func() V) *Job[V] {
	c := make(chan (chan <- V))
	done := make(chan struct{})

	go work(f, c, done)

	return &Job[V]{ send: c, done: done }
}


func (p *Job[V]) put(c chan <- V) bool {
	select {
	case p.send <- c:
		return true
	case <- p.done:
		return false
	}
}

func (p *Job[V]) Wait() (V, error) {
	return p.WaitContext(context.Background())
}

func (p *Job[V]) WaitContext(ctx context.Context) (V, error) {
	c := make(chan V)
	var v V

	if !p.put(c) {
		return v, ErrSendReceiver
	}

	select {
	case v, ok := <- c:
		if ok {
			return v, nil
		}
		return v, ErrReceiverClosed
	case <- p.done:
		return v, ErrAlreadyDone
	case <- ctx.Done():
		return v, ctx.Err()
	}
}

func (p *Job[V]) Channel() <- chan V {
	c := make(chan V)

	if !p.put(c) {
		// Fail to put receive channel
		close(c)
	} else {
		go func(){
			<- p.done
			close(c)
		}()
	}

	return c
}


func First[V any](jobs ...*Job[V]) (V, error) {
	c := make(chan V)
	done := true

	for _, job := range jobs {
		if job.put(c) {
			done = false
		}
	}

	if done {
		var v V
		return v, ErrAlreadyDone
	}

	v, err := <-c
	return v, err
}


func MaybeAll[V any](jobs ...*Job[V]) []WithError[V] {
	vs := make([]WithError[V], 0, len(jobs))

	for _, job := range jobs {
		v, err := job.Wait()
		vs = append(vs, WithError[V]{ Value: v, Error: err })
	}

	return vs
}
