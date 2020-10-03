package v2

import (
	"compress/zlib"
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
	"io"
)

type Decoder struct {
	reader *Reader

	packetFactory func(state common.PacketState, id common.PacketID) common.Packet
	threshold     int
}

func NewDecoder(r io.Reader) *Decoder {
	decoder := &Decoder{}

	decoder.reader = NewReader(r)

	return decoder
}

func (d *Decoder) SetThreshold(threshold int) {
	d.threshold = threshold
}

func (d *Decoder) SetPacketFactory(factory func(state common.PacketState, id common.PacketID) common.Packet) {
	d.packetFactory = factory
}

func (d Decoder) Unmarshal(state common.PacketState) (common.Packet, error) {
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

	packet := d.packetFactory(state, common.PacketID(packetId))
	if packet == nil {
		packet = &BasePacket{
			id: common.PacketID(packetId),
		}
	}

	err = packet.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed reading packet payload: %w", err)
	}

	return packet, nil
}

func (d Decoder) decompressBuffer(buf common.Buffer) (common.Buffer, error) {
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
