package json

import (
	"github.com/ymd-h/go/encoding/json"
)

type (
	Encoder struct {
		json.Encoder
	}
	Decoder = json.Decoder
)


func (_ Encoder) ContentType() string {
	return "application/json"
}
