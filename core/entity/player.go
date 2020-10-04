package entity

import (
	"github.com/paulhobbel/performcraft/core/game"
)

type Player interface {
	LivingEntity
	Sender

	IsOnline() bool
	SetOnline(state bool)

	GetProfile() *game.Profile
}
