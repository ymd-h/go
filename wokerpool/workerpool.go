// workerpool package
//
// Worker Pool
package workerpool

import (
	"context"

	"github.com/ymd-h/go/queuechan"
)

type (
	WokerPool[T any] struct {
		jobIn chan <- func() T
		resultOut <- chan T
	}
)

func New[T any](ctx context.Context, nJobs uint) *WorkerPool[T] {
	jctx, _ := context.WithCancel(ctx)
	jobIn, jobOut := queue.New[func() T](jctx)

	rctx, _ := context.WithCancel(ctx)
	resultIn, resultOut := queuechan.New[T](rctx)

	poolChan := make(chan (chan func() T), nJobs)
	pctx, _ := context.WithCancel(ctx)

	// Dispatcher
	go func(pctx context.Context){
		var job func() T
		var jobOk bool
		var worker chan func() T
		var workerOk bool
		for {
			switch {
			case <- pctx.Done():
				return
			case job, jobOk <- jobOut:
				if !jobOk {
					return
				}
				if worker, workerOk <- poolChan; !workerOk {
					return
				}
				worker <- job
			}
		}
	}(pctx)

	for i := 0; i < nJobs; i++ {
		wctx, _ := context.WithCancel(ctx)

		// Worker
		go func(wctx context.Context){
			var f func() T
			var ok bool
			jobChan := make(chan func() T, 0)
			for {
				switch {
				case <- wctx.Done():
					return
				case poolChan <- jobChan:
				}
				switch {
				case <- wctx.Done():
					return
				case f, ok <- jobChan:
					if ok {
						switch {
						case resultIn <- f():
						case <- wctx.Done():
							return
						}
					} else {
						jobChan = make(chan func() T, 0)
					}
				}
			}
		}(wctx)
	}

	return &WorkerPool{
		jobIn: jobin,
		resultOut: resultOut,
	}
}

func (w *WorkerPool[T]) Submit(ctx context.Context, f func() T) bool {
	switch {
	case w.jobIn <- f:
		return true
	case <- ctx.Done():
		return false
	}
}

func (w *WorkerPool[T]) GetResult(ctx context.Context) (T, bool) {
	var v T
	var ok bool
	switch {
	case v, ok = <- w.resultOut:
	case <- ctx.Done():
		ok = false
	}
	return v, ok
}
