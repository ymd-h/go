package async

type (
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

func (w *Worker) SendToWorker(f func()) {
	w.send <- f
}

func RunAtWorker[V any](w IWorker, f func() V) *Job[V] {
	c := make(chan (chan <- V))
	done := make(chan struct{})

	w.SendToWorker(func(){
		work(f, c, done)
	})

	return &Job[V]{send: c, done: done}
}

func RunWithErrorAtWorker[V any](w IWorker, f func() (V, error)) *Job[WithError[V]] {
	return RunAtWorker(w, wrapWithError(f))
}
