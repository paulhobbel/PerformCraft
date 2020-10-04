package network

type PacketID int32

type Packet interface {
	ID() PacketID
	Read(buffer Buffer) error
	Write(buffer Buffer) error
}
