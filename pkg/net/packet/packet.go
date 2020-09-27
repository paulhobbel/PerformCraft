package packet

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

type Packet struct {
	ID   int32
	Data []byte
}

func Marshal(id int32, fields ...FieldEncoder) (pk Packet) {
	pk.ID = id

	for _, field := range fields {
		pk.Data = append(pk.Data, field.Encode()...)
	}

	return
}

func (p Packet) Scan(fields ...FieldDecoder) error {
	r := bytes.NewReader(p.Data)
	for _, v := range fields {
		err := v.Decode(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet) Pack(threshold int) (pack []byte) {
	data := append(VarInt(p.ID).Encode(), p.Data...)

	if threshold > 0 {
		if len(data) > threshold {
			length := VarInt(len(data)).Encode()
			data = compress(data)

			pack = append(pack, VarInt(len(length)+len(data)).Encode()...)
			pack = append(pack, length...)
			pack = append(pack, data...)
		} else {
			pack = append(pack, VarInt(len(data)+1).Encode()...)
			pack = append(pack, 0x00)
			pack = append(pack, data...)
		}

		return
	}

	pack = append(pack, VarInt(int32(len(data))).Encode()...)
	pack = append(pack, data...)

	return
}

func ReadPacket(r Reader, compressed bool) (*Packet, error) {
	var length VarInt
	if err := length.Decode(r); err != nil {
		return nil, err
	}

	if length < 1 {
		return nil, fmt.Errorf("packet lenght too short")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("failed reading data of packet: %v", err)
	}

	if compressed {
		return decompress(data)
	}

	buf := bytes.NewBuffer(data)
	var packetID VarInt
	if err := packetID.Decode(buf); err != nil {
		return nil, fmt.Errorf("failed reading packet id: %v", err)
	}

	return &Packet{
		ID:   int32(packetID),
		Data: buf.Bytes(),
	}, nil
}

func compress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()

	return b.Bytes()
}

func decompress(data []byte) (*Packet, error) {
	r := bytes.NewReader(data)

	var size VarInt
	if err := size.Decode(r); err != nil {
		return nil, err
	}

	uncompressedData := make([]byte, size)
	if size != 0 {
		zr, err := zlib.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("decompress failed: %v", err)
		}
		defer zr.Close()

		_, err = io.ReadFull(zr, uncompressedData)
		if err != nil {
			return nil, fmt.Errorf("decompress failed: %v", err)
		}
	} else {
		uncompressedData = data[1:]
	}

	buf := bytes.NewBuffer(uncompressedData)
	var packetID VarInt
	if err := packetID.Decode(buf); err != nil {
		return nil, fmt.Errorf("failed reading packet id: %v", err)
	}

	return &Packet{
		ID:   int32(packetID),
		Data: buf.Bytes(),
	}, nil

}
