package model

import "github.com/gofrs/uuid/v5"

type Playlist struct {
	ID          uuid.UUID
	Name        string
	UserID      uuid.UUID
	TrackIDList []uuid.UUID
}
