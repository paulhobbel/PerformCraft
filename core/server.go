package core

import (
	"github.com/paulhobbel/performcraft/core/entity"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/network"
)

type Server interface {
	GetNetwork() network.Network

	GetPlayers() []entity.Player
	GetVersion() game.Version
}
