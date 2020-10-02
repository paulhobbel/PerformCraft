package v2

import (
	"bufio"
	"bytes"
	"io"
)

type Decoder struct {
	buf       *Buffer
	threshold int
}

func NewDecoder(r io.Reader) *Decoder {
	decoder := &Decoder{}

	if dr, ok := r.(*bytes.Buffer); ok {
		decoder.buf = &Buffer{dr}
	} else {
		decoder.buf = &Buffer{bytes.Buffer.ReadFrom()}
	}

	return decoder
}

func (d Decoder) SetThreshold(threshold int) {
	d.threshold = threshold
}

func (d Decoder) Unmarshal(v interface{}) {

}
