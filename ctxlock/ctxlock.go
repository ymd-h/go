// Package ctxlock provides context aware lock
package ctxlock

import (
	"context"
	"sync/atomic"
)

type (
	// Lock implements ordinary exclusive lock.
	Lock struct {
		write chan struct{}
	}

	// SharableLock implements exclusive lock for writer and shared lock for reader.
	SharableLock struct {
		lock *Lock
		want *atomic.Int32
		add chan struct{}
		done chan struct{}
	}

	UnlockFunc func()
)

// onceFunc returns wrapped function which can execute only once.
// This is simpler reimplementation of sync.Once.Do,
// because we don't need to wait unlock function.
func onceFunc(f func()) func() {
	var done atomic.Bool

	return func(){
		if done.CompareAndSwap(false, true) {
			f()
		}
	}
}


// NewLock creates a new Lock and returns the pointer to it.
func NewLock() *Lock {
	return &Lock{
		write: make(chan struct{}, 1),
	}
}

// Lock tries to lock and returns unlock function when it succeed.
// If ctx is canceled, lock is canceled and context.Cause(ctx) error is returned.
func (L *Lock) Lock(ctx context.Context) (UnlockFunc, error) {
	// If ctx has already been canceled, we don't try to lock at all.
	select {
	case <- ctx.Done():
		return nil, context.Cause(ctx)
	default:
	}

	select {
	case L.write <- struct{}{}:
		return L.unlockFunc(), nil
	case <- ctx.Done():
		return nil, context.Cause(ctx)
	}
}


// unlockFunc returns unlock function.
// It is safe to call the returned function multiple time.
func (L *Lock) unlockFunc() UnlockFunc {
	return onceFunc(func(){ <-L.write })
}

// NewSharableLock creates a new SharableLock and returns the pointer to it.
func NewSharableLock() *SharableLock {
	var want atomic.Int32
	return &SharableLock{
		lock: NewLock(),
		want: &want,
		add: make(chan struct{}),
		done: make(chan struct{}),
	}
}

// readThread takes unlock function for already locked exclusive lock
// and tracks the number of reader locks.
// Once all the readers are finished, unlock is called.
func (L *SharableLock) readThread(unlock UnlockFunc){
	defer unlock()

	i := 1
	for i > 0 {
		if L.want.Load() > 0 {
			_, ok := <- L.done
			if !ok {
				panic("BUG: done channel should not be closed.")
			}
			i -= 1
		} else {
			select {
			case _, ok := <- L.add:
				if !ok {
					panic("BUG: add channel should not be closed.")
				}
				i += 1
			case _, ok := <- L.done:
				if !ok {
					panic("BUG: done channel should not be closed.")
				}
				i -= 1
			}
		}
	}
}

// doneFunc returns done function for reader.
// It is safe to call the returned function multiple time.
func (L *SharableLock) doneFunc() UnlockFunc {
	return onceFunc(func(){ L.done <- struct{}{} })
}

// SharedLock tries to lock for reader and returns unlock function when it succeed.
// If ctx is canceled, lock is canceled and context.Cause(ctx) error is returned.
func (L *SharableLock) SharedLock(ctx context.Context) (UnlockFunc, error) {
	select {
	case <- ctx.Done():
		return nil, context.Cause(ctx)
	default:
	}

	select {
	case L.add <- struct{}{}:
	case L.lock.write <- struct{}{}:
		go L.readThread(L.lock.unlockFunc())
	case <- ctx.Done():
		return nil, context.Cause(ctx)
	}

	return L.doneFunc(), nil
}

// ExclusiveLock tries to lock for writer and returns unlock function when it succeed.
// If ctx is canceled, lock is canceled and context.Cause(ctx) error is returned.
func (L *SharableLock) ExclusiveLock(ctx context.Context) (UnlockFunc, error) {
	select {
	case <- ctx.Done():
		return nil, context.Cause(ctx)
	default:
	}

	L.want.Add(1)
	defer L.want.Add(-1)

	return L.lock.Lock(ctx)
}

// UnlockOnCancel schedules to call unlock when ctx cancels.
// If the ctx has already been canceled, unlock is called immediately.
func (f UnlockFunc) UnlockOnCancel(ctx context.Context){
	context.AfterFunc(ctx, f)
}
