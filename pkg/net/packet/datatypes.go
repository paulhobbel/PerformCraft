package packet

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"math"
)

// Boolean of true is encoded as 0x01, false as 0x00.
type Boolean bool

// Encode a Boolean.
func (b Boolean) Encode() []byte {
	if b {
		return []byte{0x01}
	}

	return []byte{0x00}
}

// Decode a Boolean.
func (b *Boolean) Decode(r Reader) error {
	v, err := r.ReadByte()
	if err != nil {
		return err
	}

	*b = v != 0
	return nil
}

//Byte is signed 8-bit integer, two's complement
type Byte int8

// Encode a Byte.
func (b Byte) Encode() []byte {
	return []byte{byte(b)}
}

// Decode a Byte.
func (b *Byte) Decode(r Reader) error {
	v, err := r.ReadByte()
	if err != nil {
		return err
	}

	*b = Byte(v)
	return nil
}

// UnsignedByte is an unsigned 8-bit integer.
type UnsignedByte uint8

// Encode a UnsignedByte.
func (b UnsignedByte) Encode() []byte {
	return []byte{byte(b)}
}

// Decode a UnsignedByte.
func (b *UnsignedByte) Decode(r Reader) error {
	v, err := r.ReadByte()
	if err != nil {
		return err
	}

	*b = UnsignedByte(v)
	return nil
}

// Short is a signed 16-bit integer, two's complement.
type Short int16

// Encode a Short.
func (s Short) Encode() []byte {
	n := uint16(s)
	return []byte{
		byte(n >> 8),
		byte(n),
	}
}

// Decode a Short.
func (s *Short) Decode(r Reader) error {
	v, err := readNBytes(r, 2)
	if err != nil {
		return err
	}

	*s = Short(int16(v[0])<<8 | int16(v[1]))
	return nil
}

// UnsignedShort is an unsigned 16-bit integer.
type UnsignedShort uint16

// Encode a UnsignedShort.
func (s UnsignedShort) Encode() []byte {
	n := uint16(s)
	return []byte{
		byte(n >> 8),
		byte(n),
	}
}

// Decode a UnsignedShort.
func (s *UnsignedShort) Decode(r Reader) error {
	v, err := readNBytes(r, 2)
	if err != nil {
		return err
	}

	*s = UnsignedShort(int16(v[0])<<8 | int16(v[1]))
	return nil
}

// Int is a signed 32-bit integer, two's complement.
type Int int32

// Encode a Int.
func (s Int) Encode() []byte {
	n := uint32(s)
	return []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}

// Decode a Int.
func (s *Int) Decode(r Reader) error {
	v, err := readNBytes(r, 4)
	if err != nil {
		return err
	}

	*s = Int(int32(v[0])<<24 | int32(v[1])<<16 | int32(v[2])<<8 | int32(v[3]))
	return nil
}

// Long is a signed 64-bit integer, two's complement.
type Long int64

// Encode a Long.
func (s Long) Encode() []byte {
	n := uint64(s)
	return []byte{
		byte(n >> 56), byte(n >> 48), byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
	}
}

// Decode a Long.
func (s *Long) Decode(r Reader) error {
	v, err := readNBytes(r, 8)
	if err != nil {
		return err
	}

	*s = Long(
		int64(v[0])<<56 | int64(v[1])<<48 | int64(v[2])<<40 | int64(v[3])<<32 |
			int64(v[4])<<24 | int64(v[5])<<16 | int64(v[6])<<8 | int64(v[6]))
	return nil
}

// Float is a single-precision 32-bit IEE 753 floating point number.
type Float float32

// Encode a Float.
func (f Float) Encode() []byte {
	return Int(math.Float32bits(float32(f))).Encode()
}

// Decode a Float.
func (f *Float) Decode(r Reader) error {
	var v Int
	if err := v.Decode(r); err != nil {
		return err
	}

	*f = Float(math.Float32frombits(uint32(v)))
	return nil
}

// Double is a double-precision 64-bit IEE 753 floating point number.
type Double float64

// Encode a Double.
func (d Double) Encode() []byte {
	return Long(math.Float64bits(float64(d))).Encode()
}

// Decode a Double.
func (d *Double) Decode(r Reader) error {
	var v Long
	if err := v.Decode(r); err != nil {
		return err
	}

	*d = Double(math.Float64frombits(uint64(v)))
	return nil
}

// String is a sequence of unicode scalar values.
type String string

// Encode a String.
func (s String) Encode() (r []byte) {
	bytes := []byte(s)
	r = append(r, VarInt(len(bytes)).Encode()...)
	r = append(r, bytes...)

	return
}

// Decode a String.
func (s *String) Decode(r Reader) error {
	var length VarInt
	if err := length.Decode(r); err != nil {
		return err
	}

	bs, err := readNBytes(r, int(length))
	if err != nil {
		return err
	}

	*s = String(bs)
	return nil
}

type VarInt int32

//Encode a VarInt
func (v VarInt) Encode() (vi []byte) {
	num := uint32(v)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		vi = append(vi, byte(b))
		if num == 0 {
			break
		}
	}
	return
}

//Decode a VarInt
func (v *VarInt) Decode(r Reader) error {
	var n uint32
	for i := 0; ; i++ {
		sec, err := r.ReadByte()
		if err != nil {
			return err
		}

		n |= uint32(sec&0x7F) << uint32(7*i)

		if i >= 5 {
			return errors.New("VarInt is too big")
		} else if sec&0x80 == 0 {
			break
		}
	}

	*v = VarInt(n)
	return nil
}

type VarLong int64

// Position x as a 26-bit integer, followed by y as a 12-bit integer, followed by z as a 26-bit integer (all signed, two's complement)
type Position struct {
	X int
	Y int
	Z int
}

// Encode a Position.
func (p Position) Encode() []byte {
	result := make([]byte, 8)
	raw := uint64(p.X&0x3FFFFFF)<<38 | uint64(p.Z&0x3FFFFFF)<<12 | uint64(p.Y&0xFFF)

	for i := 7; i >= 0; i-- {
		result[i] = byte(raw)
		raw >>= 8
	}

	return result
}

func (p *Position) Decode(r Reader) error {
	var v Long
	if err := v.Decode(r); err != nil {
		return err
	}

	x := int(v >> 8)
	y := int(v & 0xFFF)
	z := int(v << 26 >> 38)

	if x >= 1<<25 {
		x -= 1 << 26
	}

	if y >= 1<<11 {
		y -= 1 << 12
	}

	if z >= 1<<25 {
		z -= 1 << 26
	}

	p.X, p.Y, p.Z = x, y, z

	return nil
}

type UUID uuid.UUID

func (u UUID) Encode() []byte {
	return u[:]
}

func (u UUID) Decode(r Reader) error {
	_, err := io.ReadFull(r, u[:])
	return err
}

func readNBytes(r Reader, n int) (bs []byte, err error) {
	bs = make([]byte, n)

	for i := 0; i < n; i++ {
		bs[i], err = r.ReadByte()
		if err != nil {
			return
		}
	}

	return
}
