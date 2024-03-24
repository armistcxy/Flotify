package model

import "github.com/gofrs/uuid/v5"

type Artist struct {
	ID          uuid.UUID
	Name        string
	Description string
}
