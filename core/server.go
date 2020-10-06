package core

import (
	"github.com/paulhobbel/performcraft/core/entity"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/network"
	"time"
)

type Server interface {
	GetNetwork() network.Network

	GetPlayers() []entity.Player
	GetVersion() game.Version

	Start()
}

type serverImpl struct {
	network network.Network
	version game.Version

	lastTick int64
}

func NewServer() Server {
	return &serverImpl{
		network: network.NewNetwork(),
		version: game.NewVersion("PerformCraft 1.15.2", 578),
	}
}

func (s serverImpl) GetNetwork() network.Network {
	return s.network
}

func (serverImpl) GetPlayers() []entity.Player {
	panic("implement me")
}

func (s serverImpl) GetVersion() game.Version {
	return s.version
}

func (s serverImpl) Start() {
	go s.tick()
	s.network.Listen(25569)
}

func (s *serverImpl) tick() {
	tick := time.NewTicker(time.Millisecond * (1000 / 20))
	defer tick.Stop()

	for {
		select {
		case curr := <-tick.C:
			unix := curr.UnixNano() / 1e6

			s.network.Tick(int(unix - s.lastTick))

			//log.Printf("[Server]: Ticking..., curr: %v, last: %v", curr.UnixNano() / 1e6, s.lastTick)
			s.lastTick = unix
		}
	}
}

var instance Server

func Instance() Server {
	if instance == nil {
		instance = NewServer()
	}

	return instance
}
