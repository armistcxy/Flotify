package handler

import (
	"context"
	"flotify/internal/model"
	"flotify/internal/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type TrackHandler struct {
	repository repository.TrackRepository
}

func NewTrackHandler(repo repository.TrackRepository) TrackHandler {
	return TrackHandler{
		repository: repo,
	}
}

func (th *TrackHandler) CreateTrack(c *gin.Context) {
	type RequestTrack struct {
		Name   string `json:"name"`
		Length int    `json:"length"`
	}

	request_track := RequestTrack{}
	if err := c.BindJSON(&request_track); err != nil {
		log.Println("JSON struct is wrong")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	track := &model.Track{
		Name:   request_track.Name,
		Length: request_track.Length,
	}

	track, err := th.repository.CreateTrack(context.Background(), track)
	if err != nil {
		log.Println("create function is wrong")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, track)
}

func (th *TrackHandler) GetTrackByID(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	track, err := th.repository.GetTrackByID(context.Background(), id)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, track)
}
