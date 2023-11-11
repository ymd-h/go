package async

import (
	"fmt"
	"testing"
	"time"
)


func TestRun(t *testing.T) {
	tests := []struct{
		name string
		f func() int
		want int
	}{
		{
			name: "soon",
			f: func() int { return 1 },
			want: 1,
		},
		{
			name: "wait",
			f: func() int {
				_ = <- time.After(time.Duration(1e+3))
				return 2
			},
			want: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T){
			job := Run(test.f)

			v, err := job.Wait()
			if err != nil {
				t.Errorf("Fail to Wait\n")
				return
			}
			if test.want != v {
				t.Errorf("Fail: %d != %d\n", v, test.want)
				return
			}

			_, err = job.Wait()
			if err == nil {
				t.Errorf("Must not success to Wait\n")
				return
			}

			ch := job.Channel()
			select {
			case _, ok := <- ch:
				if ok {
					t.Errorf("Closed channel must fail to extract\n")
					return
				}
			default:
				t.Errorf("Done Job must return closed channel\n")
				return
			}
		})
	}

}


func TestFirst(t *testing.T){
	n := 5
	jobs := make([]*Job[struct{}], 0, n)

	for i := 0; i < n; i++ {
		jobs = append(jobs, Run(func () struct{} {
			return struct{}{}
		}))
	}

	success := 0
	for {
		if _, err := First(jobs...); err != nil {
			break
		}
		success += 1
	}
	if success != n {
		t.Errorf("All Job must success: %d != %d\n", n, success)
		return
	}
}


func TestFirstDivide(t *testing.T){
	n := 5
	jobs := make([]*Job[struct{}], 0, n)

	for i := 0; i < n; i++ {
		jobs = append(jobs, Run(func() struct{} {
			return struct{}{}
		}))
	}

	success := 0
	_, err := First(jobs...)
	if err == nil {
		success += 1
	}

	for {
		if _, err = First(jobs[:n/2]...); err != nil {
			break
		}
		success += 1
	}

	for {
		if _, err = First(jobs[n/2:]...); err != nil {
			break
		}
		success += 1
	}

	if success != n {
		t.Errorf("All jobs must success: %d != %d\n", n, success)
		return
	}
}



func TestRunError(t *testing.T){
	t.Run("error", func(_ *testing.T){
		job := Run(WrapErrorFunc(func() (struct{}, error){
			return struct{}{}, fmt.Errorf("Error")
		}))

		v, err := job.Wait()
		if err != nil {
			t.Errorf("Must Success\n")
			return
		}

		if v.Error == nil {
			t.Errorf("Must have error\n")
			return
		}
	})
}
