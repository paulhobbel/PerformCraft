package network

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/common"
	"io"
)

type Network interface {
	GetSession(id uuid.UUID) (Session, bool)

	Listen(port int) error

	common.Tickable
	io.Closer
}
