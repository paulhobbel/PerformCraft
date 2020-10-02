package common

import "fmt"

type PacketID int32
type PacketState int

const (
	Handshaking PacketState = iota
	Status
	Login
	Play
)

func (s PacketState) String() string {
	switch s {
	case Handshaking:
		return "Handshaking"
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
