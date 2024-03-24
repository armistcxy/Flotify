package model

import "github.com/gofrs/uuid/v5"

type Artist struct {
	ID          uuid.UUID `example:"3983a1d6-759b-4e5e-b307-7b7e06a05a85"`
	Name        string
	Description string
}
