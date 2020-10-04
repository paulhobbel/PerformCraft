package entity

import "github.com/google/uuid"

type Entity interface {
	UUID() uuid.UUID
}
