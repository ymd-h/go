package qchan

import (
	"context"
	"testing"
	"time"
)


func TestQueue(t *testing.T){
	q := New[int]()

	for i := 0; i < 5; i ++ {
		select {
		case q.In() <- i:
		case <- q.Done():
			t.Errorf("Fail: %d\n", i)
			return
		}
	}

	LOOP:
	for i := 0; i < 5; i++ {
		select {
		case v, ok := <- q.Out():
			if ok && (v == i) {
				continue LOOP
			}
		case <- q.Done():
		}
		t.Errorf("Fail: %d\n", i)
		return
	}

	select {
	case <- q.Out():
		t.Errorf("Fail\n")
		return
	case <- time.After(time.Duration(1000)):
	}

	close(q.In())
	<- q.Done()
	_, ok := <- q.Out()
	if ok {
		t.Errorf("Fail\n")
		return
	}
}


func TestQueueContext(t *testing.T){
	ctx, cancel := context.WithCancel(context.Background())
	q := NewWithContext[int](ctx)

	for i := 0; i < 5; i++ {
		select {
		case q.In() <- i:
		case <- q.Done():
			t.Errorf("Fail: %d\n", i)
			return
		}
	}

	LOOP:
	for i := 0; i < 5; i++ {
		select {
		case v, ok := <- q.Out():
			if ok && (v == i) {
				continue LOOP
			}
		case <- q.Done():
		}
	}

	cancel()
	<- q.Done()

	select {
	case q.In() <- -1:
		t.Errorf("Fail\n")
		return
	case <- time.After(time.Duration(5000)):
	}

	_, ok := <- q.Out()
	if ok {
		t.Errorf("Fail\n")
		return
	}

	if q.Error() == nil {
		t.Errorf("Must Error\n")
		return
	}
}
