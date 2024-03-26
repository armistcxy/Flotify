package model

import "github.com/gofrs/uuid/v5"

type Artist struct {
	ID          uuid.UUID `example:"3983a1d6-759b-4e5e-b307-7b7e06a05a85"`
	Name        string    `example:"Taylor Swift"`
	Description string    `example:"Taylor Swift (born December 13, 1989, West Reading, Pennsylvania, U.S.) is a multitalented singer-songwriter and global superstar who has captivated audiences with her heartfelt lyrics and catchy melodies, solidifying herself as one of the most influential artists in contemporary music."`
}

type Artists struct {
	Artists []Artist `swaggertype:"object,string" example:"key:value"`
}
