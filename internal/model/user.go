package model

import "github.com/gofrs/uuid/v5"

type User struct {
	ID uuid.UUID
	Name string
}