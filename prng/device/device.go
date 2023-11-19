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
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return binary.NativeEndian.Uint32(b)
}


func (_ Device64) Next() uint64 {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return binary.NativeEndian.Uint64(b)
}
