package device

import (
	"crypto/rand"
	"math/big"
)

type (
	Device32 struct {}
	Device64 struct {}
)


func (_ Device32) Next() uint32 {
	m := big.NewInt(0xFF_FF_FF_FF)
	s, _ := rand.Int(rand.Reader, m)
	return uint32(s.Uint64())
}


func (_ Device64) Next() uint64 {
	m := big.NewInt(0)
	m.SetBytes([]byte{
		255, 255, 255, 255,
		255, 255, 255, 255,
	})
	s, _ := rand.Int(rand.Reader, m)
	return s.Uint64()
}
