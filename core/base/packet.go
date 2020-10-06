package base

import (
	"github.com/paulhobbel/performcraft/core/bufio"
)

type PacketID int32

type Packet interface {
	ID() PacketID
	Read(buffer bufio.ByteBuffer) error
	Write(buffer bufio.ByteBuffer) error
}
