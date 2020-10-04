package impl

import (
	"github.com/paulhobbel/performcraft/core"
	"github.com/paulhobbel/performcraft/core/entity"
	"github.com/paulhobbel/performcraft/core/game"
	base_network "github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/impl/network"
)

type serverImpl struct {
	network base_network.Network
}

func NewServer() core.Server {
	return &serverImpl{
		network: network.NewNetwork(),
	}
}

func (s serverImpl) GetNetwork() base_network.Network {
	return s.network
}

func (serverImpl) GetPlayers() []entity.Player {
	panic("implement me")
}

func (serverImpl) GetVersion() game.Version {
	panic("implement me")
}
