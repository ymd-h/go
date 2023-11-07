package async

import (
	"time"
)


type (
	WithError[V any] struct {
		Value V
		Error error
	}

	Optional[V any] struct {
		Value V
		OK bool
	}

	Job[V any] struct {
		send chan <- (chan <- V)
		done <- chan struct{}
	}
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

func (p *Job[V]) GetWait() (V, bool) {
	c := make(chan V)

	if p.put(c) {
		select {
		case v, ok := <- c:
			return v, ok
		case <- p.done:
		}

	}

	var v V
	return v, false
}

func (p *Job[V]) GetWaitDuration(d time.Duration) (V, bool) {
	c := make(chan V)
	if p.put(c) {
		select {
		case v, ok := <- c:
			return v, ok
		case <- p.done:
		case <- time.After(d):
		}
	}

	var v V
	return v, false
}

func (p *Job[V]) Get() (V, bool) {
	c := make(chan V)
	if p.put(c) {
		select {
		case v, ok := <- c:
			return v, ok
		case <- p.done:
		default:
		}
	}

	var v V
	return v, false
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


func First[V any](jobs ...*Job[V]) (V, bool) {
	c := make(chan V)
	done := true

	for _, job := range jobs {
		if job.put(c) {
			done = false
		}
	}

	if done {
		var v V
		return v, false
	}

	v, ok := <-c
	return v, ok
}


func MaybeAll[V any](jobs ...*Job[V]) []Optional[V] {
	vs := make([]Optional[V], 0, len(jobs))

	for _, job := range jobs {
		v, ok := vs.GetWait()
		vs = append(vs, Optional[V]{ Value: v, OK: ok })
	}

	return vs
}
