package common

import "github.com/google/uuid"

type Player struct {
	UUID uuid.UUID
	Name string
}
