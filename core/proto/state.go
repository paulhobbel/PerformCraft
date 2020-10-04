package proto

import "fmt"

type ProtocolState int

const (
	Handshake ProtocolState = iota
	Status
	Login
	Play
	Unknown
)

func (s ProtocolState) String() string {
	switch s {
	case Handshake:
		return "Handshake"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Play:
		return "Play"
	default:
		panic(fmt.Errorf("no state for value %d", s))
	}
}

func ProtocolStateFromId(id int32) ProtocolState {
	switch id {
	case 0:
		return Handshake
	case 1:
		return Status
	case 2:
		return Login
	case 3:
		return Play
	default:
		return Unknown
	}
}
