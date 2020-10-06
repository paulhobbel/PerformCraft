package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/paulhobbel/performcraft/core/nbt"
	"github.com/paulhobbel/performcraft/core/text"
	"io"
	"math"
)

type ByteReader interface {
	io.Reader
	io.ByteReader

	ReadBoolean() (bool, error)

	ReadShort() (int16, error)

	ReadUShort() (uint16, error)

	ReadInt() (int32, error)

	ReadLong() (int64, error)

	ReadFloat() (float32, error)

	ReadDouble() (float64, error)

	ReadString() (string, error)

	ReadMessage() (text.Message, error)

	ReadVarInt() (int32, error)

	ReadNbt(interface{}) error
}

type byteReader interface {
	io.Reader
	io.ByteReader
}

type readerImpl struct {
	parent byteReader
}

func NewByteReader(r io.Reader) ByteReader {
	reader := &readerImpl{}

	if br, ok := r.(byteReader); ok {
		reader.parent = br
	} else {
		reader.parent = bufio.NewReader(r)
	}

	return reader
}

func (r readerImpl) Read(v []byte) (int, error) {
	return r.parent.Read(v)
}

func (r readerImpl) ReadByte() (byte, error) {
	return r.parent.ReadByte()
}

func (r readerImpl) ReadBoolean() (bool, error) {
	v, err := r.ReadByte()
	if err != nil {
		return false, err
	}

	return v != 0, nil
}

func (r readerImpl) ReadShort() (int16, error) {
	v, err := r.ReadUShort()
	return int16(v), err
}

func (r readerImpl) ReadUShort() (uint16, error) {
	bs, err := r.readNBytes(2)
	if err != nil {
		return 0, err
	}
	return uint16(bs[0])<<8 | uint16(bs[1]), nil
}

func (r readerImpl) ReadInt() (int32, error) {
	v, err := r.readNBytes(4)
	if err != nil {
		return 0, err
	}

	return int32(v[0])<<24 | int32(v[1])<<16 | int32(v[2])<<8 | int32(v[3]), nil
}

func (r readerImpl) ReadLong() (int64, error) {
	v, err := r.readNBytes(8)
	if err != nil {
		return 0, err
	}

	return int64(v[0])<<56 | int64(v[1])<<48 | int64(v[2])<<40 | int64(v[3])<<32 |
		int64(v[4])<<24 | int64(v[5])<<16 | int64(v[6])<<8 | int64(v[6]), nil
}

func (r readerImpl) ReadFloat() (float32, error) {
	v, err := r.ReadInt()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(uint32(v)), nil
}

func (r readerImpl) ReadDouble() (float64, error) {
	v, err := r.ReadLong()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(uint64(v)), nil
}

func (r readerImpl) ReadString() (string, error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}

	bs, err := r.readNBytes(int(length))
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (r readerImpl) ReadMessage() (msg text.Message, err error) {
	length, err := r.ReadVarInt()
	if err != nil {
		return text.Message{}, err
	}

	err = json.NewDecoder(io.LimitReader(r, int64(length))).Decode(&msg)

	return
}

func (r readerImpl) ReadVarInt() (int32, error) {
	var v int32

	for i := 0; ; i++ {
		tmp, err := r.ReadByte()
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

func (r *readerImpl) ReadNbt(v interface{}) error {
	return nbt.NewDecoder(r.parent).Unmarshal(v)
}

func (r readerImpl) readNBytes(n int) (bs []byte, err error) {
	bs = make([]byte, n)

	for i := 0; i < n; i++ {
		bs[i], err = r.ReadByte()
		if err != nil {
			return
		}
	}

	return
}
