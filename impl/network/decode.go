package network

import (
	"compress/zlib"
	"fmt"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
	"github.com/paulhobbel/performcraft/pkg/common"
	"io"
)

type packetDecoder struct {
	reader network.Reader

	packetFactory func(state proto.ProtocolState, id network.PacketID) network.Packet
	threshold     int
}

func NewDecoder(r io.Reader) *packetDecoder {
	decoder := &packetDecoder{}

	decoder.reader = NewReader(r)

	return decoder
}

func (d *packetDecoder) SetThreshold(threshold int) {
	d.threshold = threshold
}

func (d *packetDecoder) SetPacketFactory(factory func(state proto.ProtocolState, id network.PacketID) network.Packet) {
	d.packetFactory = factory
}

func (d packetDecoder) Decode(state proto.ProtocolState) (network.Packet, error) {
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
	buf := NewBufferFrom(data)

	// Decompress if needed
	buf, err = d.decompressBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("decompress failed: %w", err)
	}

	packetId, err := buf.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed reading packet id: %w", err)
	}

	packet := d.packetFactory(state, network.PacketID(packetId))
	if packet == nil {
		packet = &network.BasePacket{
			id: network.PacketID(packetId),
		}
	}

	err = packet.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed reading packet payload: %w", err)
	}

	return packet, nil
}

func (d packetDecoder) decompressBuffer(buf network.Buffer) (network.Buffer, error) {
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

			return NewBufferFrom(decompressed), nil
		}
	}

	return buf, nil
}
