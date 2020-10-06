package nbt

import (
	"bufio"
	"io"
)

const (
	TagEnd byte = iota
	TagByte
	TagShort
	TagInt
	TagLong
	TagFloat
	TagDouble
	TagByteArray
	TagString
	TagList
	TagCompound
	TagIntArray
	TagLongArray
	TagNone = 0xFF
)

type DecoderReader interface {
	io.ByteReader
	io.Reader

	Peek(n int) ([]byte, error)
}

type Decoder struct {
	reader DecoderReader
}

func NewDecoder(r io.Reader) *Decoder {
	decoder := &Decoder{}

	if dr, ok := r.(DecoderReader); ok {
		decoder.reader = dr
	} else {
		decoder.reader = bufio.NewReader(r)
	}

	return decoder
}
