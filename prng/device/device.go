package device

import (
	"crypto/rand"
	"encoding/binary"
)

type (
	Device32 struct {}
	Device64 struct {}
)


func (_ Device32) Next() uint32 {
	m := big.NewInt(0x01_00_00_00_00)
	s, _ := rand.Int(rand.Reader, m)
	return uint32(s.Uint64())
}


func (_ Device64) Next() uint64 {
	m := big.NewInt(0)
	m.SetBytes([]byte{
		1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	})
	s, _ := rand.Int(rand.Reader, m)
	return s.Uint64()
}
