package async

import (
	"context"
	"testing"
	"time"
)



func TestWorker(t *testing.T){
	w := NewWorker(context.Background(), 1)

	job, err := RunAtWorker(context.Background(), w, func() int { return 1 })
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	v, err := job.Wait()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	if v != 1 {
		t.Errorf("%d != 1\n", v)
		return
	}

	job, err = RunAtWorker(context.Background(), w, func() int { return 2 })
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	ctxT, _ := context.WithTimeout(context.Background(), time.Duration(1 * 1000))
	_, err = RunAtWorker(ctxT, w, func() int { return 3 })
	if err == nil {
		t.Errorf("Must Fail\n")
		return
	}
}


func TestLazyWorker(t *testing.T){
	w := NewLazyWorker(context.Background(), 2)

	job, err := RunAtWorker(context.Background(), w, func() int { return 2 })
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	v, err := job.Wait()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	if v != 2 {
		t.Errorf("%d != 2\n", v)
		return
	}

	jobs := make([]*Job[int], 0, 2)
	for i := 0; i < 2; i++ {
		job, err = RunAtWorker(context.Background(), w, func() int {
			return i
		})
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		jobs = append(jobs, job)
	}

	ctxT, _ := context.WithTimeout(context.Background(), time.Duration(1 * 10000))
	_, err = RunAtWorker(ctxT, w, func() int {
		return 5
	})
	if err == nil {
		t.Errorf("Must Fail\n")
		return
	}

	for _, job := range jobs {
		_, err = job.Wait()
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
	}

	ctxT, _ = context.WithTimeout(context.Background(), time.Duration(1 * 10000))
	job, err = RunAtWorker(ctxT, w, func() int { return -1 })
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	v, err = job.Wait()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
	if v != -1 {
		t.Errorf("%d != -1\n", v)
		return
	}
}
