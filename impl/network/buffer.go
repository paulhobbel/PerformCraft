package network

import (
	"bytes"
	"errors"
	"github.com/paulhobbel/performcraft/core/network"
	"math"
)

type bufferImpl struct {
	network.Reader

	buf *bytes.Buffer
}

func NewEmptyBuffer() network.Buffer {
	buf := &bytes.Buffer{}

	return &bufferImpl{
		Reader: NewReader(buf),
		buf:    buf,
	}
}

func NewBufferFrom(data []byte) network.Buffer {
	buf := bytes.NewBuffer(data)

	return &bufferImpl{
		Reader: NewReader(buf),
		buf:    buf,
	}
}

func (b bufferImpl) Len() int {
	return b.buf.Len()
}

func (b bufferImpl) Cap() int {
	return b.buf.Cap()
}

func (b bufferImpl) Bytes() []byte {
	return b.buf.Bytes()
}

func (b bufferImpl) Write(v []byte) (int, error) {
	return b.buf.Write(v)
}

func (b bufferImpl) WriteByte(v byte) error {
	return b.buf.WriteByte(v)
}

func (b bufferImpl) WriteBoolean(v bool) error {
	if v {
		return b.WriteByte(0x01)
	}
	return b.WriteByte(0x00)
}

func (b bufferImpl) WriteShort(v int16) error {
	return b.WriteUShort(uint16(v))
}

func (b *bufferImpl) WriteUShort(v uint16) error {
	_, err := b.Write([]byte{
		byte(v >> 8),
		byte(v),
	})
	return err
}

func (b bufferImpl) WriteInt(v int32) error {
	n := uint32(v)

	_, err := b.Write([]byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})

	return err
}

func (b bufferImpl) WriteLong(v int64) error {
	n := uint64(v)

	_, err := b.Write([]byte{
		byte(n >> 56), byte(n >> 48), byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
	})

	return err
}

func (b bufferImpl) WriteFloat(v float32) error {
	return b.WriteInt(int32(math.Float32bits(v)))
}

func (b bufferImpl) WriteDouble(v float64) error {
	return b.WriteLong(int64(math.Float64bits(v)))
}

func (b bufferImpl) WriteString(v string) (err error) {
	bs := []byte(v)
	err = b.WriteVarInt(int32(len(bs)))
	_, err = b.Write(bs)

	return
}

//func (b bufferImpl) WriteMessage(v chat.Message) error {
//	bs, err := json.Marshal(v)
//	if err != nil {
//		return err
//	}
//
//	return b.WriteString(string(bs))
//}

func (b bufferImpl) WriteVarInt(v int32) (err error) {
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

func (b bufferImpl) WriteNbt(v interface{}) error {
	return errors.New("not implemented")
}
