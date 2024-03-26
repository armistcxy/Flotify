package model

import "github.com/gofrs/uuid/v5"

type Track struct {
	ID       uuid.UUID   `example:"3983a1d6-759b-4e5e-b307-7b7e06a05a85"`
	Name     string      `example:"Blue Town"`
	Length   int         `example:"88"`
	ArtistID []uuid.UUID `example:"3983a1d6-759b-4e5e-b307-7b7e06a05a85"`
}

type Tracks struct {
	Tracks []Track `swaggertype:"object,string" example:"key:value"`
}
