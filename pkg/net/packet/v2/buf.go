package v2

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/nbt"
	"math"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(b *bytes.Buffer) *Buffer {
	return &Buffer{b}
}

func (b Buffer) WriteBoolean(v bool) error {
	if v {
		return b.WriteByte(0x01)
	}
	return b.WriteByte(0x00)
}

func (b Buffer) ReadBoolean() (bool, error) {
	v, err := b.ReadByte()
	if err != nil {
		return false, err
	}

	return v != 0, nil
}

func (b Buffer) WriteShort(v int16) error {
	return b.WriteUShort(uint16(v))
}

func (b Buffer) ReadShort() (int16, error) {
	v, err := b.ReadUShort()
	return int16(v), err
}

func (b *Buffer) WriteUShort(v uint16) error {
	_, err := b.Write([]byte{
		byte(v >> 8),
		byte(v),
	})
	return err
}

func (b Buffer) ReadUShort() (uint16, error) {
	bs, err := b.readNBytes(2)
	if err != nil {
		return 0, err
	}
	return uint16(bs[0])<<8 | uint16(bs[1]), nil
}

func (b Buffer) WriteInt(v int32) error {
	n := uint32(v)

	_, err := b.Write([]byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})

	return err
}

func (b Buffer) ReadInt() (int32, error) {
	v, err := b.readNBytes(4)
	if err != nil {
		return 0, err
	}

	return int32(v[0])<<24 | int32(v[1])<<16 | int32(v[2])<<8 | int32(v[3]), nil
}

func (b Buffer) WriteLong(v int64) error {
	n := uint64(v)

	_, err := b.Write([]byte{
		byte(n >> 56), byte(n >> 48), byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
	})

	return err
}

func (b Buffer) ReadLong() (int64, error) {
	v, err := b.readNBytes(8)
	if err != nil {
		return 0, err
	}

	return int64(v[0])<<56 | int64(v[1])<<48 | int64(v[2])<<40 | int64(v[3])<<32 |
		int64(v[4])<<24 | int64(v[5])<<16 | int64(v[6])<<8 | int64(v[6]), nil
}

func (b Buffer) WriteFloat(v float32) error {
	return b.WriteInt(int32(math.Float32bits(v)))
}

func (b Buffer) ReadFloat() (float32, error) {
	v, err := b.ReadInt()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(uint32(v)), nil
}

func (b Buffer) WriteDouble(v float64) error {
	return b.WriteLong(int64(math.Float64bits(v)))
}

func (b Buffer) ReadDouble() (float64, error) {
	v, err := b.ReadLong()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(uint64(v)), nil
}

func (b Buffer) WriteString(v string) (err error) {
	bs := []byte(v)
	err = b.WriteVarInt(int32(len(bs)))
	_, err = b.Write(bs)

	return
}

func (b Buffer) ReadString() (string, error) {
	length, err := b.ReadVarInt()
	if err != nil {
		return "", err
	}

	bs, err := b.readNBytes(int(length))
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (b Buffer) WriteVarInt(v int32) (err error) {
	for {
		tmp := v & 0x7F
		v >>= 7

		if v != 0 {
			tmp |= 0x80
		}

		err = b.WriteByte(byte(tmp))
		if err != nil || v == 0 {
			break
		}
	}

	return
}

func (b Buffer) ReadVarInt() (int32, error) {
	var v int32

	for i := 0; ; i++ {
		tmp, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		v |= int32(tmp&0x7F) << uint(i*7)

		if i >= 5 {
			return 0, fmt.Errorf("VarInt is too big, %d > 5", i)
		} else if tmp&0x80 == 0 {
			break
		}
	}

	return v, nil
}

func (b Buffer) WriteNbt(v interface{}) error {
	return errors.New("not implemented")
}

func (b *Buffer) ReadNbt(v interface{}) error {
	return nbt.NewDecoder(b).Unmarshal(v)
}

func (b Buffer) readNBytes(n int) (bs []byte, err error) {
	bs = make([]byte, n)

	for i := 0; i < n; i++ {
		bs[i], err = b.ReadByte()
		if err != nil {
			return
		}
	}

	return
}
