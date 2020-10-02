package v2

type Packet interface {
	Read(buffer Buffer) error
	Write(buffer Buffer) error
}
