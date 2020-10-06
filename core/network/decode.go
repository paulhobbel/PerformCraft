package network

import (
	"compress/zlib"
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
	"github.com/paulhobbel/performcraft/core/proto"
	"io"
)

type packetDecoder struct {
	reader bufio.ByteReader

	protocolState base.ProtocolState
	threshold     int
}

func NewPacketDecoder(r io.Reader) *packetDecoder {
	return &packetDecoder{
		reader:        bufio.NewByteReader(r),
		protocolState: base.Handshake,
		threshold:     0,
	}
}

func (d *packetDecoder) SetProtocolState(state base.ProtocolState) {
	d.protocolState = state
}

func (d *packetDecoder) SetThreshold(threshold int) {
	d.threshold = threshold
}

func (d packetDecoder) Decode() (base.Packet, error) {
	length, err := d.reader.ReadVarInt()
	if err != nil {
		return nil, err
	}

	// Read full packet and create a buffer
	data := make([]byte, length)
	_, err = io.ReadFull(d.reader, data)
	if err != nil {
		return nil, fmt.Errorf("failed reading data of packet: %w", err)
	}
	buf := bufio.NewByteBufferFrom(data)

	// Decompress if needed
	buf, err = d.decompressBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("decompress failed: %w", err)
	}

	packetId, err := buf.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed reading packet id: %w", err)
	}

	packet := proto.GetClientPacket(d.protocolState, base.PacketID(packetId))
	if packet == nil {
		return nil, fmt.Errorf("failed reading packet payload: unknown packet id")
	}

	err = packet.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed reading packet payload: %w", err)
	}

	return packet, nil
}

func (d packetDecoder) decompressBuffer(buf bufio.ByteBuffer) (bufio.ByteBuffer, error) {
	if d.threshold > 0 {
		size, err := buf.ReadVarInt()
		if err != nil {
			return buf, err
		}

		if size != 0 {
			zReader, err := zlib.NewReader(buf)
			if err != nil {
				return buf, err
			}
			defer zReader.Close()

			// Allocate raw buffer
			decompressed := make([]byte, size)
			_, err = io.ReadFull(zReader, decompressed)
			if err != nil {
				return buf, err
			}

			return bufio.NewByteBufferFrom(decompressed), nil
		}
	}

	return buf, nil
}
