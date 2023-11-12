package async

import (
	"context"
	"testing"
	"time"
)

func TestWorker(t *testing.T){
	tests := []struct{
		name string
		w IWorker
	}{
		{
			name: "Worker",
			w: NewWorker(context.Background(), 1),
		},
		{
			name: "LazyWorker",
			w: NewLazyWorker(context.Background(), 1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(*testing.T){
			job, err := RunAtWorker(
				context.Background(),
				test.w,
				func() struct{} { return struct{}{} },
			)
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			_, err = job.Wait()
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			ctx, cancel := context.WithCancel(context.Background())

			_, err = RunAtWorker(
				context.Background(),
				test.w,
				func() struct{} {
					<-ctx.Done()
					return struct{}{}
				},
			)
			if err != nil {
				t.Errorf("Fail: %v\n", err)
				return
			}

			ctxT, _ := context.WithTimeout(context.Background(), time.Duration(1000))
			_, errT := RunAtWorker(
				ctxT,
				test.w,
				func() struct{} { return struct{}{} },
			)
			if errT == nil {
				t.Errorf("Must Fail")
				return
			}

			cancel()

			_, err = RunAtWorker(
				context.Background(),
				test.w,
				func() struct{} { return struct{}{} },
			)
		})
	}
}
