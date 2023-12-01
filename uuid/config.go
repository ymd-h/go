package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type (
	ITimestamp interface {
		UnixMilli() int64
	}

	defaultTimestamp struct {}

	timestampFunc struct {
		f func() int64
	}

	IRandom interface {
		FillRandom([]byte) error
	}

	defaultRandom struct {}

	randomFunc struct {
		f func([]byte) error
	}

	Config struct {
		t ITimestamp
		r IRandom
	}
)

var (
	ErrTimestampAlreadySet = errors.New("Timestamp has already been set.")
	ErrRandomAlreadySet = errors.New("Random has already been set.")
	ErrTimestampOutOfRange = errors.New("Timestamp is out of range.")
	ErrUnknownOption = errors.New("Unknown config option")
)


func (_ defaultTimestamp) UnixMilli() int64 {
	return time.Now().UnixMilli()
}

func TimestampFrom(f func() int64) *timestampFunc {
	return &timestampFunc{ f: f }
}

func (t *timestampFunc) UnixMilli() int64 {
	return t.f()
}

func (_ defaultRandom) FillRandom(b []byte) error {
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Errorf("Fail default FillRandom(): %v", err)
	}
	return nil
}

func RandomFrom(f func([]byte) error) *randomFunc {
	return &randomFunc{ f: f }
}

func (r *randomFunc) FillRandom(b []byte) error {
	return r.f(b)
}

func NewConfig(options ...any) (*Config, error) {
	c := Config{
		t: nil,
		r: nil,
	}

	for _, o := range options {
		if t, ok := o.(ITimestamp); ok {
			if c.t != nil {
				return nil, ErrTimestampAlreadySet
			}

			c.t = t
		} else if r, ok := o.(IRandom); ok {
			if c.r != nil {
				return nil, ErrRandomAlreadySet
			}
			c.r = r
		} else {
			return nil, fmt.Errorf("%w: %T", ErrUnknownOption, o)
		}
	}

	if c.t == nil {
		c.t = defaultTimestamp{}
	}
	if c.r == nil {
		c.r = defaultRandom{}
	}

	return &c, nil
}

func (c *Config) UUIDv4() (*UUIDv4, error) {
	var u UUIDv4
	err := c.r.FillRandom(u.b[:])
	if err != nil {
		return nil, fmt.Errorf("UUIDv4: Fail to Fill Random: %w", err)
	}

	u.setVersion(4)
	u.setVariant(uuidVariant, uuidVariantMask)

	return &u, nil
}

func (c *Config) UUIDv7() (*UUIDv7, error) {
	var u UUIDv7

	unix_ms := c.t.UnixMilli()
	if (unix_ms > 0xFFFFFFFFFFFF) || (unix_ms < 0) {
		return nil, fmt.Errorf("%w: %v", ErrTimestampOutOfRange, unix_ms)
	}

	binary.BigEndian.PutUint16(u.b[ :2], uint16(unix_ms >> 32))
	binary.BigEndian.PutUint32(u.b[2:6], uint32(unix_ms & 0xFFFFFFFF))

	err := c.r.FillRandom(u.b[6:])
	if err != nil {
		return nil, fmt.Errorf("UUIDv7: Fail to Fill Random: %w", err)
	}

	u.setVersion(7)
	u.setVariant(uuidVariant, uuidVariantMask)

	return &u, nil
}

func NewUUIDv4() (*UUIDv4, error) {
	cfg, _ := NewConfig()
	return cfg.UUIDv4()
}

func NewUUIDv7() (*UUIDv7, error) {
	cfg, _ := NewConfig()
	return cfg.UUIDv7()
}
