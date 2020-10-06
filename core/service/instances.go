package service

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/game"
)

var gameProfileServiceInst GameProfileService

func GetGameProfileService() GameProfileService {
	if gameProfileServiceInst == nil {
		gameProfileServiceInst = &gameProfileServiceImpl{
			profiles: make(map[uuid.UUID]*game.Profile),
		}
	}

	return gameProfileServiceInst
}
