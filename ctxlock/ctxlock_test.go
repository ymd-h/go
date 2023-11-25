package ctxlock

import (
	"context"
	"testing"
	"time"
)

func newTimeout(dt time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dt)
}


func TestLock(t *testing.T){
	L := NewLock()
	dt := time.Duration(10000)

	// Lock
	// -> OK
	ctx, _ := newTimeout(dt)
	unlock, err := L.Lock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	// Lock when it is already locked.
	// -> error
	ctxT, _ := newTimeout(dt)
	_, errT := L.Lock(ctxT)
	if errT == nil {
		t.Errorf("Must Fail\n")
		return
	}

	unlock()

	// Lock after unlock
	// -> OK
	ctx, _ = newTimeout(dt)
	unlock, err = L.Lock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	unlock()

	// Lock with already canceled context
	ctx, cancel := newTimeout(dt)
	cancel()
	_, errC := L.Lock(ctx)
	if errC == nil {
		t.Errorf("Must Fail\n")
		return
	}


	// Lock when it unlock during waiting.
	// -> OK
	unlock, err = L.Lock(context.Background())
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	go func(){
		<- time.After(time.Duration(100000000))
		unlock()
	}()

	_, err = L.Lock(context.Background())
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
}


func TestSharableLock(t *testing.T){
	L := NewSharableLock()
	dt := time.Duration(100000)

	// ExclusiveLock
	// -> OK
	ctx, _ := newTimeout(dt)
	unlock, err := L.ExclusiveLock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	// SharedLock when it has already been ExclusiveLock-ed.
	// -> error
	ctxT, _ := newTimeout(dt)
	_, errT := L.SharedLock(ctxT)
	if errT == nil {
		t.Errorf("Must Fail\n")
		return
	}

	unlock()

	// SharedLock after unlock
	// -> OK
	ctx, _ = newTimeout(dt)
	unlock1, err := L.SharedLock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	// SharedLock when it is already been SharedLock-ed.
	// -> OK
	ctx, _ = newTimeout(dt)
	unlock2, err := L.SharedLock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	// ExclusiveLock when it is already been SharedLock-ed.
	ctx, _ = newTimeout(dt)
	_, errE := L.ExclusiveLock(ctx)
	if errE == nil {
		t.Errorf("Must Fail\n")
		return
	}

	unlock1()
	unlock2()
}


func TestUnlock(t *testing.T){
	L := NewLock()
	dt := time.Duration(100000)

	ctx, _ := newTimeout(dt)
	unlock, err := L.Lock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	// Call multiple time safely.
	unlock()
	unlock()

	ctx, _ = newTimeout(dt)
	unlock, err = L.Lock(ctx)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	unlock()
}
