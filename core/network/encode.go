package network

import (
	"github.com/paulhobbel/performcraft/core/base"
	buf2 "github.com/paulhobbel/performcraft/core/bufio"
	"io"
)

type packetEncoder struct {
	writer    io.Writer
	threshold int
}

func NewPacketEncoder(w io.Writer) *packetEncoder {
	return &packetEncoder{w, 0}
}

func (e *packetEncoder) SetThreshold(threshold int) {
	e.threshold = threshold
}

func (e packetEncoder) Encode(p base.Packet) (err error) {
	buf := buf2.NewByteBuffer()
	packetBuf := buf2.NewByteBuffer()

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
