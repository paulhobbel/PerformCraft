package proto

import (
	"fmt"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type ServerPacketKeepAlive struct {
	KeepAliveID int64
}

func (ServerPacketKeepAlive) ID() base.PacketID {
	return PlayKeepAliveServer
}

func (p *ServerPacketKeepAlive) Read(buffer bufio.ByteBuffer) (err error) {
	p.KeepAliveID, err = buffer.ReadLong()

	return
}

func (p ServerPacketKeepAlive) Write(buffer bufio.ByteBuffer) error {
	return buffer.WriteLong(p.KeepAliveID)
}

type ServerPacketPlayJoinGame struct {
	EntityId          int32
	GameMode          byte
	Dimension         int32
	Seed              int64
	MaxPlayers        byte
	LevelType         string
	ViewDistance      int32
	ReducedDebugInfo  bool
	EnableSpawnScreen bool
}

func (ServerPacketPlayJoinGame) ID() base.PacketID {
	return PlayJoinGame
}

func (p *ServerPacketPlayJoinGame) Read(b bufio.ByteBuffer) (err error) {
	p.EntityId, err = b.ReadInt()
	p.GameMode, err = b.ReadByte()
	p.Dimension, err = b.ReadInt()
	p.Seed, err = b.ReadLong()
	p.MaxPlayers, err = b.ReadByte()
	p.LevelType, err = b.ReadString()
	p.ViewDistance, err = b.ReadVarInt()
	p.ReducedDebugInfo, err = b.ReadBoolean()
	p.EnableSpawnScreen, err = b.ReadBoolean()

	return
}

func (p ServerPacketPlayJoinGame) Write(b bufio.ByteBuffer) (err error) {
	err = b.WriteInt(p.EntityId)
	err = b.WriteByte(p.GameMode)
	err = b.WriteInt(p.Dimension)
	err = b.WriteLong(p.Seed)
	err = b.WriteByte(p.MaxPlayers)
	err = b.WriteString(p.LevelType)
	err = b.WriteVarInt(p.ViewDistance)
	err = b.WriteBoolean(p.ReducedDebugInfo)
	err = b.WriteBoolean(p.EnableSpawnScreen)

	return
}

type ServerPacketPlayPlayerPositionLook struct {
	X          float64
	Y          float64
	Z          float64
	Yaw        float32
	Pitch      float32
	Flags      byte
	TeleportID int32
}

func (ServerPacketPlayPlayerPositionLook) ID() base.PacketID {
	return PlayPlayerPositionLookServer
}

func (p *ServerPacketPlayPlayerPositionLook) Read(b bufio.ByteBuffer) (err error) {
	p.X, err = b.ReadDouble()
	p.Y, err = b.ReadDouble()
	p.Z, err = b.ReadDouble()
	p.Yaw, err = b.ReadFloat()
	p.Pitch, err = b.ReadFloat()
	p.Flags, err = b.ReadByte()
	p.TeleportID, err = b.ReadVarInt()

	return
}

func (p ServerPacketPlayPlayerPositionLook) Write(b bufio.ByteBuffer) (err error) {
	err = b.WriteDouble(p.X)
	err = b.WriteDouble(p.Y)
	err = b.WriteDouble(p.Z)
	err = b.WriteFloat(p.Yaw)
	err = b.WriteFloat(p.Pitch)
	err = b.WriteByte(p.Flags)
	err = b.WriteVarInt(p.TeleportID)

	return
}

func (p ServerPacketPlayPlayerPositionLook) String() string {
	return fmt.Sprintf("ServerPacketPlayPlayerPositionLook{X: %v, Y: %v, Z: %v, Yaw: %v, Pitch: %v, Flags: %v, TeleportID: %v}",
		p.X, p.Y, p.Z, p.Yaw, p.Pitch, p.Flags, p.TeleportID)
}
