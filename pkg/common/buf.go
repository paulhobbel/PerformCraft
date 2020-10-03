package common

import "github.com/paulhobbel/performcraft/pkg/net/chat"

type Buffer interface {
	Len() int
	Cap() int
	Bytes() []byte

	Write(v []byte) (int, error)
	Read(v []byte) (int, error)

	WriteByte(v byte) error
	ReadByte() (byte, error)

	WriteBoolean(v bool) error
	ReadBoolean() (bool, error)

	WriteShort(v int16) error
	ReadShort() (int16, error)

	WriteUShort(v uint16) error
	ReadUShort() (uint16, error)

	WriteInt(v int32) error
	ReadInt() (int32, error)

	WriteLong(v int64) error
	ReadLong() (int64, error)

	WriteFloat(v float32) error
	ReadFloat() (float32, error)

	WriteDouble(v float64) error
	ReadDouble() (float64, error)

	WriteString(v string) error
	ReadString() (string, error)

	WriteMessage(v chat.Message) error
	ReadMessage() (chat.Message, error)

	WriteVarInt(v int32) error
	ReadVarInt() (int32, error)

	WriteNbt(v interface{}) error
	ReadNbt(interface{}) error
}
