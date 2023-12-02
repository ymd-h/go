package prng

import (
	"math"
)

type (
	iRandom interface { uint32  | uint64  }
	fRandom interface { float32 | float64 }

	IPRNG[I iRandom] interface {
		Next() I
	}

	prng32on64 struct {
		g IPRNG[uint64]
	}

	IRandom[I iRandom, F fRandom] interface {
		Uint() I
		Float() F
	}

	random32 struct {
		g IPRNG[uint32]
	}

	random64 struct {
		g IPRNG[uint64]
	}

	ISplittable[T any] interface {
		Copy() T
		Jump()
	}
)

func (p *prng32on64) Next() uint32 {
	return uint32(p.g.Next() >> 32)
}

func NewRandom32[I iRandom](g IPRNG[I]) IRandom[uint32, float32] {
	if g32, ok := g.(IPRNG[uint32]); ok {
		return &random32{g: &pring32on64{g: g32}}
	}

	if g64, ok := g.(IPRNG[uint64]); ok {
		return &random32{g: g64}
	}

	panic("BUG")
}

func NewRandom64(g IPRNG[uint64]) IRandom[uint64, float64] {
	return &random64{g: g}
}

func (p *random32) Uint() uint32 {
	return p.g.Next()
}

func (p *random32) Float() float32 {
	x := p.g.Next()
	f := math.Float32frombits((uint32(0x7F) << 23) | (x >> 9))
	return f - 1.0
}

func (p *random64) Uint() uint64 {
	return p.g.Next()
}

func (p *random64) Float() float64 {
	x := p.g.Next()
	f := math.Float64frombits((uint64(0x3FF) << 52) | (x >> 12))
	return f - 1.0
}

func Split[T ISplittable[T]](r T) T {
	c := r.Copy()
	c.Jump()
	return c
}
