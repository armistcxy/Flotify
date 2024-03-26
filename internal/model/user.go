package model

import "github.com/gofrs/uuid/v5"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}
