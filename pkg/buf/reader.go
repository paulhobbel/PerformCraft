package buf

import "io"

type Reader interface {
	io.ByteReader
	io.Reader

	Peek(n int) ([]byte, error)
}
