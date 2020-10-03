package v2

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	"io"
)

type Encoder struct {
	writer    io.Writer
	threshold int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w, 0}
}

func (e *Encoder) SetThreshold(threshold int) {
	e.threshold = threshold
}

func (e Encoder) Marshal(p common.Packet) (err error) {
	buf := NewEmptyBuffer()
	packetBuf := NewEmptyBuffer()

	err = packetBuf.WriteVarInt(int32(p.ID()))
	err = p.Write(packetBuf)
	if err != nil {
		return err
	}

	// Compressed
	if e.threshold > 0 {
		if packetBuf.Len() > e.threshold {
			// TODO: Get length of VarInt
		} else {
			buf.WriteVarInt(int32(packetBuf.Len() + 1))
			buf.WriteByte(0x00)
			buf.Write(packetBuf.Bytes())
		}
	} else {
		err = buf.WriteVarInt(int32(packetBuf.Len()))
		_, err = buf.Write(packetBuf.Bytes())
	}

	_, err = e.writer.Write(buf.Bytes())

	return
}
