package network

import (
	"bytes"
	"compress/zlib"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
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

func (e packetEncoder) Encode(p base.Packet) error {
	buf := bufio.NewByteBuffer()
	packetBuf := bufio.NewByteBuffer()

	packetBuf.WriteVarInt(int32(p.ID()))
	p.Write(packetBuf)

	// Compressed
	if e.threshold > 0 {
		deflateBuf := bufio.NewByteBuffer()
		e.deflate(packetBuf, deflateBuf)

		buf.WriteVarInt(int32(deflateBuf.Len()))
		buf.Write(deflateBuf.Bytes())
	} else {
		buf.WriteVarInt(int32(packetBuf.Len()))
		buf.Write(packetBuf.Bytes())
	}

	_, err := e.writer.Write(buf.Bytes())
	return err
}

func (e packetEncoder) deflate(src, dst bufio.ByteBuffer) {
	if src.Len() < e.threshold {
		dst.WriteByte(0x00)
		dst.Write(src.Bytes())
	} else {
		var buf bytes.Buffer
		deflater := zlib.NewWriter(&buf)
		deflater.Write(src.Bytes())
		deflater.Close()

		dst.WriteVarInt(int32(src.Len()))
		dst.Write(buf.Bytes())
	}
}
