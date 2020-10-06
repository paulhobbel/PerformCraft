package network

import "io"

type Reader interface {
	io.Reader
	io.ByteReader

	ReadBoolean() (bool, error)

	ReadShort() (int16, error)

	ReadUShort() (uint16, error)

	ReadInt() (int32, error)

	ReadLong() (int64, error)

	ReadFloat() (float32, error)

	ReadDouble() (float64, error)

	ReadString() (string, error)

	//ReadMessage() (chat.Message, error)

	ReadVarInt() (int32, error)

	ReadNbt(interface{}) error
}
