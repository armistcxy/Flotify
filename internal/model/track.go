package model

import "github.com/gofrs/uuid/v5"

type Track struct {
	ID     uuid.UUID
	Name   string
	Length int
}
