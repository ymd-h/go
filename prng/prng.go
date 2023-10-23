package prng

import (
	"crypt/rand"
	"encoding/binary"
	"math"
	"math/big"
)

type (
	IPRNG32 interface {
		next() uint32
	}

	PRNG32 struct {
		g IPRNG32
	}

	IPRNG64 interface {
		next() uint64
	}

	PRNG64 struct {
		g IPRNG64
	}
)


func Seed() uint64 {
	m := big.NewInt(0xFF_FF_FF_FF_FF_FF_FF_FF)
	s := rand.Int(rand.Reader, m)
	return s.Uint64()
}


func (p *PRNG32) Uint32() uint32 {
	return p.g.next()
}

func (p *PRNG32) Float() float32 {
	x := p.g.next()
	f := math.Float32frombits((uint32(0x7F) << 23) | (x >> 9))
	return f - 1.0
}

func (p *PRNG64) Uint64() uint64 {
	return p.g.next()
}

func (p *PRNG64) Double() float64 {
	x := p.g.next()
	f := math.Float64frombits((uint64(0x3FF) << 52) | (x >> 12))
	return f - 1.0
}
