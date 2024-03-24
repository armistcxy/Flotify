package handler

import (
	"context"
	"flotify/internal/helper"
	"flotify/internal/model"
	"flotify/internal/repository"
	"net/http"
	"strings"

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

// CreateTrack godoc
//
//	@Summary		Create a track
//	@Description	Create a new track
//	@Tags			track
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Track
//	@Failure		400	"Bad request"
//	@Router			/track [post]
func (th *TrackHandler) CreateTrack(c *gin.Context) {
	type RequestTrack struct {
		Name      string      `json:"name"`
		Length    int         `json:"length"`
		Artist_id []uuid.UUID `json:"artist_id"`
	}

	request_track := RequestTrack{}
	if err := c.BindJSON(&request_track); err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	track := &model.Track{
		Name:     request_track.Name,
		Length:   request_track.Length,
		ArtistID: request_track.Artist_id,
	}

	track, err := th.repository.CreateTrack(context.Background(), track)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, track)
}

// GetTrack godoc
//
//	@Summary		Get track's information
//	@Description	Get track's information by its ID
//	@Tags			track
//	@Param			id path string true "Track ID" example("3983a1d6-759b-4e5e-b307-7b7e06a05a85")
//	@Produce		json
//	@Success		200	{object}	model.Track
//	@Failure		400	"Bad request"
//	@Failure		500	"Internal server error"
//	@Router			/track/{id} [get]
func (th *TrackHandler) GetTrackByID(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	track, err := th.repository.GetTrackByID(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	artist_id, err := th.repository.GetArtistOfTrack(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	track.ArtistID = artist_id
	c.JSON(http.StatusOK, track)
}

func (th *TrackHandler) GetTrackWithFilter(c *gin.Context) {

	name := c.Query("name")

	page, err := helper.GetPage(c)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	limit, err := helper.GetLimit(c)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	sort_criterias_string_form := c.Query("sort")
	var sort_criterias []string
	if sort_criterias_string_form != "" {
		sort_criterias = strings.Split(sort_criterias_string_form, ",")
	}

	filter := repository.Filter{
		Props:  map[string]any{"name": name},
		Page:   page,
		Limit:  limit,
		SortBy: sort_criterias,
	}

	tracks, err := th.repository.GetTracksWithFilter(context.Background(), filter)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (th *TrackHandler) DeleteTrack(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	err = th.repository.DeleteTrack(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete successfully"})
}

func (th *TrackHandler) UpdateTrack(c *gin.Context) {
	track := model.Track{}
	if err := c.BindJSON(&track); err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if err := th.repository.UpdateTrack(context.Background(), &track); err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusAccepted, track)
}
