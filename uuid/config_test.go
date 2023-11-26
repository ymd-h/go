package uuid

import (
	"errors"
	"testing"
)


func TestNewConfig(t *testing.T) {
	_, err := NewConfig()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	T := TimestampFrom(func() int64 { return 0x0123456789ABCDEF })
	R := RandomFrom(func(b []byte) error {
		for i, _ := range b {
			b[i] = byte(i)
		}
		return nil
	})

	_, err = NewConfig(T)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	_, err = NewConfig(R)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	_, err = NewConfig(T, R)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	_, err = NewConfig(R, T)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	_, errT := NewConfig(T, T)
	if errT == nil {
		t.Errorf("Must Fail\n")
		return
	}
	if !errors.Is(errT, ErrTimestampAlreadySet) {
		t.Errorf("Fail: %v\n", errT)
		return
	}

	_, errR := NewConfig(R, R)
	if errR == nil {
		t.Errorf("Must Fail\n")
		return
	}
	if !errors.Is(errR, ErrRandomAlreadySet) {
		t.Errorf("Fail: %v\n", errR)
		return
	}
}

func TestUUIDv4(t *testing.T) {
	R := RandomFrom(func(b []byte) error {
		for i, _ := range b {
			b[i] = byte(i)
		}
		return nil
	})

	c, err := NewConfig(R)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	u, err := c.UUIDv4()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	err = u.validate()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
}

func TestUUIDv7(t *testing.T) {
	T := TimestampFrom(func() int64 { return 0x0123456789ABCDEF })
	R := RandomFrom(func(b []byte) error {
		for i, _ := range b {
			b[i] = byte(i)
		}
		return nil
	})

	c, err := NewConfig(T, R)
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	u, err := c.UUIDv7()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	err = u.validate()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
}


func TestNewUUIDv4(t *testing.T){
	u, err := NewUUIDv4()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	err = u.validate()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
}

func TestNewUUIDv7(t *testing.T){
	u, err := NewUUIDv7()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}

	err = u.validate()
	if err != nil {
		t.Errorf("Fail: %v\n", err)
		return
	}
}
