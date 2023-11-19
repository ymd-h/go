package prng

import (
	"math"
)

type (
	IPRNG32 interface {
		Next() uint32
	}

	IPRNG64 interface {
		Next() uint64
	}

	prng32on64 struct {
		g IPRNG64
	}

	IRandom32 interface {
		Uint32() uint32
		Float32() float32
	}

	IRandom64 interface {
		Uint64() uint64
		Float64() float64
	}

	random32 struct {
		g IPRNG32
	}

	random64 struct {
		g IPRNG64
	}

	ISplittable[T any] interface {
		Copy() T
		Jump()
	}
)

func Low32(g IPRNG64) *prng32on64 {
	return &prng32on64{g: g}
}

func (p *prng32on64) Next() uint32 {
	return uint32(p.Next() >> 32)
}


func NewRandom32(g IPRNG32) *random32 {
	return &random32{g: g}
}

func NewRandom64(g IPRNG64) *random64 {
	return &random64{g: g}
}

func (p *random32) Uint32() uint32 {
	return p.g.Next()
}

func (p *random32) Float32() float32 {
	x := p.g.Next()
	f := math.Float32frombits((uint32(0x7F) << 23) | (x >> 9))
	return f - 1.0
}

func (p *random64) Uint64() uint64 {
	return p.g.Next()
}

func (p *random64) Float64() float64 {
	x := p.g.Next()
	f := math.Float64frombits((uint64(0x3FF) << 52) | (x >> 12))
	return f - 1.0
}

func Split[T ISplittable[T]](r T) T {
	c := r.Copy()
	c.Jump()
	return c
}
