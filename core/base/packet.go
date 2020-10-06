package base

import "github.com/paulhobbel/performcraft/core/util"

type PacketID int32

type Packet interface {
	ID() PacketID
	Read(buffer util.ByteBuffer) error
	Write(buffer util.ByteBuffer) error
}
