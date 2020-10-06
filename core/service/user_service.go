package service

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/game"
)

type GameProfileService interface {
	GetProfileByName(name string) *game.Profile
}

type gameProfileServiceImpl struct {
	profiles map[uuid.UUID]*game.Profile
}

func (s *gameProfileServiceImpl) GetProfileByName(name string) *game.Profile {
	for _, profile := range s.profiles {
		if profile.Name == name {
			return profile
		}
	}

	// Failed to find, create new profile
	profile := &game.Profile{
		UUID: uuid.New(),
		Name: name,
	}

	s.profiles[profile.UUID] = profile

	return profile
}
