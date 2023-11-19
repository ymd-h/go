package uuid

import (
	"fmt"
)

type (
	iTextUUID interface {
		validate() error
		unmarshalText([]byte) error
	}

	iBinaryUUID interface {
		validate() error
		unmarshalBinary([]byte) error
	}

	iUUID interface {
		validate() error
		unmarshal([]byte) error
	}

	tUUID struct {
		uuid iTextUUID
	}

	bUUID struct {
		uuid iBinaryUUID
	}
)

func (t *tUUID) validate() error {
	return t.uuid.validate()
}

func (t *tUUID) unmarshal(text []byte) error {
	return t.uuid.unmarshalText(text)
}

func (b *bUUID) validate() error {
	return b.uuid.validate()
}

func (b *bUUID) unmarshal(data []byte) error {
	return b.uuid.unmarshalBinary(data)
}

func unmarshal(u iUUID, b []byte) error {
	err := u.unmarshal(b)
	if err != nil {
		return fmt.Errorf("Fail to Unmarshal: %w", err)
	}

	err = u.validate()
	if err != nil {
		return fmt.Errorf("Fail to Validate: %w", err)
	}

	return nil
}

func h2b(b byte) (byte, error) {
	switch  {
	case (0x30 <= b) && (b <= 0x39):
		//   0 -    9
		//0x30 - 0x39
		return b - 0x30, nil
	case (0x41 <= b) && (b <= 0x46):
		//    A -    F
		// 0x41 - 0x46
		return b - 0x41 + 0xA, nil
	case (0x61 <= b) && (b <= 0x66):
		//    a -    f
		// 0x61 - 0x66
		return b - 0x61 + 0xA, nil
	default:
		return 0xFF, fmt.Errorf("Invalid byte: %v", b)
	}
}

func decodeHEX(b []byte) (byte, error) {
	bu, err := h2b(b[0])
	if err != nil {
		return 0, err
	}

	bl, err := h2b(b[1])
	if err != nil {
		return 0, err
	}

	return (bu << 4) | bl, nil
}

func isHyphen(b byte) bool {
	return b == 0x2d
}
