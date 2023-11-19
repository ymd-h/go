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
	s := rand.Int(rand.Reader, m)
	return s.Uint32()
}


func (_ Device64) Next() uint64 {
	m := big.NewInt(0xFF_FF_FF_FF_FF_FF_FF_FF)
	s := rand.Int(rand.Reader, m)
	return s.Uint64()
}
