package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

type (
	ITimestamp interface {
		UnixMilli() int64
	}

	defaultTimestamp struct {}

	IRandom interface {
		FillRandom([]byte) error
	}

	defaultRandom struct {}

	Config struct {
		t ITimestamp
		r IRandom
	}
)


func (_ defaultTimestamp) UnixMilli() int64 {
	return time.Now().UnixMilli()
}

func (_ defaultRandom) FillRandom(b []byte) error {
	_, err := rand.Read(b)
	return err
}

func NewConfig(options ...any) (*Config, error) {
	c := Config{
		t: nil,
		r: nil,
	}

	for _, o := range options {
		if t, ok := o.(ITimestamp); ok {
			c.t = t
		} else if r, ok := o.(IRandom); ok {
			c.r = r
		} else {
			return nil, fmt.Errorf("Unknown Option: %T", o)
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

	u.setVersion(0b0100)
	u.setVariant(0b10)

	return &u, nil
}

func (c *Config) UUIDv7() (*UUIDv7, error) {
	var u UUIDv7

	unix_ms := c.t.UnixMilli()
	binary.BigEndian.PutUint16(u.b[ :2], uint16(unix_ms >> 32))
	binary.BigEndian.PutUint32(u.b[2:6], uint32(unix_ms & 0xFFFFFFFF))

	err := c.r.FillRandom(u.b[6:])
	if err != nil {
		return nil, fmt.Errorf("UUIDv7: Fail to Fill Random: %w", err)
	}

	u.setVersion(0b0100)
	u.setVariant(0b10)

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
