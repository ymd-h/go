package uuid

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/ymd-h/go/prng/device"
)

type(
	UUID struct {
		b [16]byte
	}

	UUIDv4 UUID
	UUIDv7 UUID
)


func FromString(s string) (*UUID, error) {
	var u UUID
	p := &u

	err := p.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (u *UUID) String() string {
	return fmt.Sprintf(
		"%02x%02x%02x%02x-" +
			"%02x%02x-" +
			"%02x%02x-" +
			"%02x%02x-" +
			"%02x%02x%02x%02x%02x%02x",
		u.b[0], u.b[1], u.b[2], u.b[3],
		u.b[4], u.b[5],
		u.b[6], u.b[7],
		u.b[8], u.b[9],
		u.b[10], u.b[11], u.b[12], u.b[13], u.b[14], u.b[15],
	)
}

func (u *UUID) Bytes() []byte {
	b := make([]byte, 0, len(u.b))
	return append(b, u.b[:]...)
}

func (u *UUID) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, u.String())), nil
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

func (u *UUID) UnmarshalText(data []byte) error {
	n := len(data)
	if n > 2 {
		f := data[0]
		e := data[n-1]
		if (f==e) && ((f == 0x22) || (f == 0x27) || (f == 0x60)) {
			// 0x22 => "
			// 0x27 => '
			// 0x60 => `
			data = data[1:n-1]
		}
	}

	if len(data) != 36 {
		return fmt.Errorf("Wrong Data Length: %d", len(data))
	}

	offset := 0
	for i := 0; i < 20; i++ {
		switch i {
		case 4, 7, 10, 13:
			if !isHyphen(data[2*i-offset]) {
				return fmt.Errorf("Hypen (-) is missing at %d.", 2*i-offset)
			}
			offset += 1
		default:
			b, err := decodeHEX(data[2*i-offset:2*i-offset+2])
			if err != nil {
				return err
			}
			u.b[i-offset] = b
		}
	}

	return nil
}

func (u *UUID) MarshalBinary() ([]byte, error) {
	return u.Bytes(), nil
}

func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("UUID UnmarshalBinary: Wrong Byte Length: %d", len(data))
	}

	for i := 0; i < 16; i++ {
		u.b[i] = data[i]
	}

	return nil
}


func (u *UUID) TryUUIDv4() (*UUIDv4, error) {
	version := u.b[6] >> 4
	if version != 4 {
		return nil, fmt.Errorf("Version is not 4: %d", version)
	}

	variant := u.b[8] >> 6
	if variant != 8 {
		return nil, fmt.Errorf("Variant is not 0x10: %d", variant)
	}

	return &UUIDv4{ b: u.b }, nil
}


func NewUUIDv4() *UUIDv4 {
	prng := device.Device64{}

	r1 := prng.Next()
	r2 := prng.Next()

	var u UUIDv4
	binary.NativeEndian.PutUint64(u.b[:4], r1)
	binary.NativeEndian.PutUint64(u.b[4:], r2)

	var version byte = 0b0100
	var variant byte = 0b10
	u.b[6] = (u.b[6] & 0x0F) | (version << 4)
	u.b[8] = (u.b[8] & 0b00111111) | (variant << 6)

	return &u
}

func (u *UUIDv4) UnmarshalText(text []byte) error {
	var uuid UUID

	err := uuid.UnmarshalText(text)
	if err != nil {
		return err
	}

	u4, err := uuid.TryUUIDv4()
	if err != nil {
		return err
	}

	u.b = u4.b
	return nil
}

func (u *UUIDv4) UnmarshalBinary(data []byte) error {
	var uuid UUID

	err := uuid.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	u4, err := uuid.TryUUIDv4()
	if err != nil {
		return err
	}

	u.b = u4.b
	return nil
}

func NewUUIDv7() *UUIDv7 {
	t := time.Now()
	unix_ms := t.UnixMilli()

	p32 := device.Device32{}
	r1 := p32.Next()

	p64 := device.Device64{}
	r2 := p64.Next()

	var u UUIDv7
	binary.BigEndian.PutUint16(u.b[:2], uint16(unix_ms >> 32))
	binary.BigEndian.PutUint32(u.b[2:6], uint32(unix_ms && 0xFFFFFFFF))
	binary.NativeEndian.PutUint16(u.b[6:8], uint16(r1 >> 16))
	binary.NativeEndian.PutUint64(u.b[8:], r2)

	var version byte = 0b0100
	var variant byte = 0b10
	u.b[6] = (u.b[6] & 0x0F) | (version << 4)
	u.b[8] = (u.b[8] & 0b00111111) | (variant << 6)

	return &u
}
