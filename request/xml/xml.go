package xml

import (
	"github.com/ymd-h/go/encoding/xml"
)

type (
	Encoder struct {
		xml.Encoder
	}
	Decoder = xml.Decoder
)


func (_ Encoder) ContentType() string {
	return `application/xml: charset="UTF-8"`
}
