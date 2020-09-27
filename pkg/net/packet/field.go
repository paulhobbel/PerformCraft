package packet

import "io"

type Field interface {
	FieldEncoder
	FieldDecoder
}

type FieldEncoder interface {
	Encode() []byte
}

type FieldDecoder interface {
	Decode(r Reader) error
}

type Reader interface {
	io.Reader
	io.ByteReader
}
