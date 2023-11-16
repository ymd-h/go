package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ymd-h/go/async"
	"github.com/ymd-h/go/slices"
	"github.com/ymd-h/go/request/gob"
	"github.com/ymd-h/go/request/xml"
)

type (
	A struct {
		A1 string `json:"a1" xml:"a1"`
		A2 int `json:"a2" xml:"a2"`
		A3 B `json:"a3" xml:"a3"`
	}
	B struct {
		B1 bool `json:"b1" xml:"b1"`
		B2 []byte `json:"b2" xml:"b2"`
	}
)


func TestRequest(t *testing.T){
	http.HandleFunc("/A", func(w http.ResponseWriter, req *http.Request){
		a := A{
			A1: "a1-string",
			A2: 2,
			A3: B{
				B1: true,
				B2: []byte{0, 1, 2, 16},
			},
		}

		b, err := json.Marshal(&a)
		if err != nil {
			return
		}

		w.Write(b)
	})
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request){
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error("Fail /echo")
			return
		}
		w.Write(b)
	})
	http.HandleFunc("/ctx", func(w http.ResponseWriter, req *http.Request){
		<- req.Context().Done()
		return
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request){
		var sc int
		err := json.NewDecoder(req.Body).Decode(&sc)
		if err != nil {
			return
		}
		w.WriteHeader(sc)

		w.Write([]byte("{}"))
		return
	})

	go http.ListenAndServe("localhost:8888", nil)

	url := func(path string) string {
		return fmt.Sprintf("http://localhost:8888/%s", path)
	}

	t.Run("GET-simple", func(*testing.T){
		var a A
		resp, err := Get(url("A"), &a)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Fail StatusCode: %d\n", resp.StatusCode)
			return
		}
	})

	t.Run("GET-partial", func(*testing.T){
		var a struct {
			A1 string `json:"a1"`
			A2 int `json:"a2"`
		}
		resp, err := Get(url("A"), &a)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Fail StatusCode: %d\n", resp.StatusCode)
			return
		}
	})

	t.Run("GET-invalid-path", func(*testing.T){
		var a A
		resp, err := Get(url(""), &a)
		if err == nil {
			t.Errorf("Must Fail\n")
			return
		}
		if resp == nil {
			t.Errorf("Fail\n")
			return
		}

	})

	t.Run("POST-simple", func(*testing.T){
		a := A{
			A1: "aiu",
			A2: 25,
			A3: B{
				B1: false,
				B2: []byte{255, 255, 255},
			},
		}
		var aa A

		_, err := Post(url("echo"), &a, &aa)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}

		if a.A1 != aa.A1 {
			t.Errorf("Fail A1: %v != %v\n", a.A1, aa.A1)
			return
		}
		if a.A2 != aa.A2 {
			t.Errorf("Fail A2: %v != %v\n", a.A2, aa.A2)
			return
		}
		if a.A3.B1 != aa.A3.B1 {
			t.Errorf("Fail B1: %v != %v\n", a.A3.B1, aa.A3.B2)
			return
		}

		if !slices.NewComparableSliceFrom(a.A3.B2).Equal(
			slices.NewComparableSliceFrom(aa.A3.B2),
		) {
			t.Errorf("Fail B2: %v != %v\n", a.A3.B2, aa.A3.B2)
			return
		}
	})

	t.Run("POST-XML", func(*testing.T){
		c := NewClient(http.DefaultClient, xml.Encoder{}, xml.Decoder{})

		a := A{
			A1: "--2",
			A2: -9,
			A3: B{
				B1: true,
				B2: []byte{0x30, 0x30, 0x30},
			},
		}
		var aa A

		_, err := c.Post(url("echo"), &a, &aa)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
		}

		if a.A1 != aa.A1 {
			t.Errorf("Fail A1: %v != %v\n", a.A1, aa.A1)
			return
		}
		if a.A2 != aa.A2 {
			t.Errorf("Fail A2: %v != %v\n", a.A2, aa.A2)
			return
		}
		if a.A3.B1 != aa.A3.B1 {
			t.Errorf("Fail B1: %v != %v\n", a.A3.B1, aa.A3.B2)
			return
		}

		if !slices.NewComparableSliceFrom(a.A3.B2).Equal(
			slices.NewComparableSliceFrom(aa.A3.B2),
		) {
			t.Errorf("Fail B2: %v != %v\n", a.A3.B2, aa.A3.B2)
			return
		}
	})

	t.Run("POST-gob", func(*testing.T){
		c := NewClient(http.DefaultClient, gob.Encoder{}, gob.Decoder{})

		a := A{
			A1: "/*/*/++",
			A2: 234,
			A3: B{
				B1: true,
				B2: []byte("098-+*"),
			},
		}
		var aa A

		_, err := c.Post(url("echo"), &a, &aa)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
		}

		if a.A1 != aa.A1 {
			t.Errorf("Fail A1: %v != %v\n", a.A1, aa.A1)
			return
		}
		if a.A2 != aa.A2 {
			t.Errorf("Fail A2: %v != %v\n", a.A2, aa.A2)
			return
		}
		if a.A3.B1 != aa.A3.B1 {
			t.Errorf("Fail B1: %v != %v\n", a.A3.B1, aa.A3.B2)
			return
		}

		if !slices.NewComparableSliceFrom(a.A3.B2).Equal(
			slices.NewComparableSliceFrom(aa.A3.B2),
		) {
			t.Errorf("Fail B2: %v != %v\n", a.A3.B2, aa.A3.B2)
			return
		}
	})

	t.Run("ctx", func(*testing.T){
		ctx, cancel := context.WithCancel(context.Background())

		job := async.Run(
			async.WrapErrorFunc(
				func()(any, error){
					return GetWithContext(ctx, url("ctx"), nil)
				},
			),
		)

		select {
		case <- job.Ready():
			t.Errorf("Fail\n")
			return
		case <- time.After(time.Duration(1000)):
		}

		cancel()

		v, err := job.Wait()
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}

		if v.Error == nil {
			t.Errorf("Must Fail\n")
			return
		}
	})

	t.Run("dispatch", func(*testing.T){
		type (
			AA struct {}
			BB struct {}
			CC struct {}
		)

		d := NewResponseDispatcher(
			DispatchItem[AA]{http.StatusOK},
			DispatchItem[BB]{http.StatusNotFound},
			DispatchItem[CC]{http.StatusTeapot},
		)

		resp, err := Post(url("status"), http.StatusOK, d)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		if _, ok := resp.Body.(*AA); !ok {
			t.Errorf("Fail: %T\n", resp.Body)
			return
		}

		resp, err = Post(url("status"), http.StatusNotFound, d)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		if _, ok := resp.Body.(*BB); !ok {
			t.Errorf("Fail: %T\n", resp.Body)
			return
		}

		resp, err = Post(url("status"), http.StatusTeapot, d)
		if err != nil {
			t.Errorf("Fail: %v\n", err)
			return
		}
		if _, ok := resp.Body.(*CC); !ok {
			t.Errorf("Fail: %T\n", resp.Body)
			return
		}

		_, err = Post(url("status"), http.StatusAccepted, d)
		if err == nil {
			t.Errorf("Must Fail")
			return
		}
	})
}
