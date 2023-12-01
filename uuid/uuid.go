package uuid

import (
	"encoding/binary"
	"fmt"
)

type(
	baseUUID struct {
		b [16]byte
	}

	UUID struct {
		baseUUID
	}

	UUIDv4 struct {
		baseUUID
	}

	UUIDv7 struct {
		baseUUID
	}
)


func (u *baseUUID) String() string {
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

func (u *baseUUID) Bytes() []byte {
	b := make([]byte, 0, len(u.b))
	return append(b, u.b[:]...)
}

func (u *baseUUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *baseUUID) MarshalBinary() ([]byte, error) {
	return u.Bytes(), nil
}

func (u *baseUUID) Version() uint8 {
	return uint8(u.b[6] >> 4)
}

func (u *baseUUID) variant() uint8 {
	return uint8(u.b[8] >> 6)
}

func (u *baseUUID) setVersion(v uint8) {
	u.b[6] = (u.b[6] & 0x0F) | (byte(v) << 4)
}

func (u *baseUUID) setVariant(v uint8) {
	u.b[8] = (u.b[8] & 0b00111111) | (byte(v) << 6)
}

func (u *baseUUID) unmarshalText(data []byte) error {
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

func (u *baseUUID) unmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("UUID UnmarshalBinary: Wrong Byte Length: %d", len(data))
	}

	for i := 0; i < 16; i++ {
		u.b[i] = data[i]
	}

	return nil
}

func (u *UUID) UnmarshalText(text []byte) error {
	return u.unmarshalText(text)
}

func (u *UUID) UnmarshalBinary(data []byte) error {
	return u.unmarshalBinary(data)
}

func (u *UUID) TryUUIDv4() (*UUIDv4, error) {
	u4 := &UUIDv4{ u.baseUUID }

	err := u4.validate()
	if err != nil {
		return nil, err
	}

	return u4, nil
}

func (u *UUID) TryUUIDv7() (*UUIDv7, error) {
	u7 := &UUIDv7{ u.baseUUID }

	err := u7.validate()
	if err != nil {
		return nil, err
	}

	return u7, nil
}

func (u *UUIDv4) UnmarshalText(text []byte) error {
	return unmarshal(&tUUID{u}, text)
}

func (u *UUIDv4) UnmarshalBinary(data []byte) error {
	return unmarshal(&bUUID{u}, data)
}

func (u *UUIDv4) validate() error {
	version := u.Version()
	if version != 4 {
		return fmt.Errorf("Version is not 4: %d", version)
	}

	variant := u.variant()
	if variant != 0b10 {
		return fmt.Errorf("Variant is not 0x10: %d", variant)
	}

	return nil
}

func (u *UUIDv7) UnmarshalText(text []byte) error {
	return unmarshal(&tUUID{u}, text)
}

func (u *UUIDv7) UnmarshalBinary(data []byte) error {
	return unmarshal(&bUUID{u}, data)
}

func (u *UUIDv7) validate() error {
	version := u.Version()
	if version != 7 {
		return fmt.Errorf("Version is not 7: %d", version)
	}

	variant := u.variant()
	if variant != 0b10 {
		return fmt.Errorf("Variant is not 0x10: %d", variant)
	}

	return nil
}

func (u *UUIDv7) timestamp() uint64 {
	return binary.BigEndian.Uint64(u.b[0:]) >> 16
}

func (u *UUIDv7) TimestampBefore(other *UUIDv7) bool {
	return u.timestamp() < other.timestamp()
}

func (u *UUIDv7) TimestampAfter(other *UUIDv7) bool {
	return u.timestamp() > other.timestamp()
}

func (u *UUIDv7) TimestampEqual(other *UUIDv7) bool {
	return u.timestamp() == other.timestamp()
}

func FromString(s string) (*UUID, error) {
	var u UUID

	err := u.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}

	return &u, nil
}
