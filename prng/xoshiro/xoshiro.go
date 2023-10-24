package xoshiro

import (
	"github.com/ymd-h/go/prng/splitmix"
)

type (
	Xoshiro256pp struct {
		s [4]uint64
	}

	Xoshiro128pp struct {
		s [4]uint32
	}
)

func rotl64(x uint64, int k) uint64 {
	return (x << k) | (x >> (64 - k))
}

func rotl32(x uint32, int k) uint32 {
	return (x << k) | (x >> (32 - k))
}

func NewXoshiro256pp(seed uint64) *Xoshiro256pp {
	sm := splitmix.NewSplitMix64(seed)

	var g Xoshiro256pp
	for i := 0; i < 4; i++ {
		g.s[i] = sm.Next()
	}

	return &g
}

func (g *Xoshiro256pp) Next() uint64 {
	result := rotl64(g.s[0] + g.s[3], 23) + g.s[0]

	t := g.s[1] << 17

	g.s[2] ^= g.s[0]
	g.s[3] ^= g.s[1]
	g.s[1] ^= g.s[2]
	g.s[0] ^= g.s[3]

	g.s[2] ^= t
	g.s[3] = rotl64(g.s[3], 45)

	return result
}

func (g *Xoshiro256pp) Jump() {
	jump := []uint64{
		0x180ec6d33cfd0aba,
		0xd5a61266f0c9392c,
		0xa9582618e03fc9aa,
		0x39abdc4529b1661c,
	}

	s := []uint64{0, 0, 0, 0}

	for _, j := range jump {
		for b := 0; b < 64; b++ {
			if(j & uint64(1) << b) != 0 {
				for i := 0; i < 4; i++ {
					s[i] ^= g.s[i]
				}
				g.Next()
			}
		}
	}

	for i := 0; i < 4; i++ {
		g.s[i] = s[i]
	}
}

func (g *Xoshiro256pp) Copy() *Xoshiro256pp {
	var c Xoshiro256pp
	for i := 0; i < 4; i++ {
		c.s[i] = g.s[i]
	}
	return &c
}

func NewXoshiro128pp(seed uint32) *Xoshiro128pp {
	sm := splitmix.NewSplitMix32(seed)

	var g Xoshiro128pp
	for i := 0; i < 4; i++ {
		g.s[i] = sm.Next()
	}

	return &g
}

func (g *Xoshiro128pp) Next() uint32 {
	reslt := rotl32(g.s[0] + g.s[3], 7) + g.s[0]

	t := g.s[1] << 9

	g.s[2] ^= g.s[0]
	g.s[3] ^= g.s[1]
	g.s[1] ^= g.s[2]
	g.s[0] ^= g.s[3]

	g.s[2] ^= t
	g.s[3] = rotl32(g.s[3], 11)

	return result
}

func (g *Xoshiro128pp) Jump() {
	jump := []uint32{0x8764000b, 0xf542d2d3, 0x6fa035c3, 0x77f2db5b}

	s := []uint32{0, 0, 0, 0}

	for _, j := range jump {
		for b := 0; b < 32; b++ {
			if (j & uint32(1) << b) != 0 {
				for i := 0; i < 4; i++ {
					s[i] ^= g.s[i]
				}
			}
			g.Next()
		}
	}

	for i := 0; i < 4; i++ {
		g.s[i] = s[i]
	}
}

func (g *Xoshiro128pp) Copy() *Xoshiro128pp {
	var c Xoshiro128pp
	for i := 0; i < 4; i++ {
		c.s[i] = g.s[i]
	}
	return &c
}
