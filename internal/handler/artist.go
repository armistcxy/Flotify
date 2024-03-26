package handler

import (
	"context"
	"flotify/internal/helper"
	"flotify/internal/model"
	"flotify/internal/repository"
	"flotify/internal/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type ArtistHandler struct {
	repository repository.ArtistRepository
}

func NewArtistHandler(repo repository.ArtistRepository) ArtistHandler {
	return ArtistHandler{
		repository: repo,
	}
}

// CreateArtist godoc
//
//		@Summary		Create an artist
//		@Description	Create a new artist
//		@Tags			artists
//		@Accept			json
//		@Produce		json
//	 	@Param 			artist body model.Artist true "Artist Information"
//		@Success		200	{object}	model.Artist
//		@Failure		400
//		@Failure		500
//		@Router			/artists [post]
func (ah *ArtistHandler) CreateArtist(c *gin.Context) {
	type RequestArtist struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	request_artist := RequestArtist{}
	err := c.BindJSON(&request_artist)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
	}

	artist := &model.Artist{
		Name:        request_artist.Name,
		Description: request_artist.Description,
	}

	artist, err = ah.repository.CreateArtist(context.Background(), artist)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, artist)
}

// GetArtistInformation godoc
//
//		@Summary		Get information of an artist
//		@Description	Get information of an artist using ID
//		@Tags			artists
//		@Produce		json
//	 	@Param	 		id path string true "Artist ID" example("3983a1d6-759b-4e5e-b307-7b7e06a05a85")
//		@Success		200	{object}	model.Artist
//		@Failure		400 "Bad Request"
//		@Failure		500 "Internal Server Error"
//		@Router			/artists/{id} [get]
func (ah *ArtistHandler) GetInfoArtistByID(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	artist, err := ah.repository.GetArtistByID(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, artist)
}

//	 GetTopTracksOfArist godoc
//
//		@Summary		Get top tracks of an artist
//		@Description	Get top tracks of an artist using ID
//		@Tags			artists
//		@Param 			id path string true "Artist ID" example("3983a1d6-759b-4e5e-b307-7b7e06a05a85")
//		@Produce		json
//		@Success		200	{object}	model.Tracks
//		@Failure		400 "Bad Request"
//		@Failure		500 "Internal Server Error"
//		@Router			/artists/{id}/tracks [get]
func (ah *ArtistHandler) GetArtistTracksByID(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	tracks, err := ah.repository.GetTrackOfArtist(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	response := gin.H{
		"tracks": tracks,
	}
	c.JSON(http.StatusOK, response)
}

// GetListOfArtist godoc
//
//	@Summary		Get list of artists
//	@Description	Get list of artists that satisfied conditions in filter
//	@Tags			artists
//	@Produce		json
//	@Param 			name query string false "name of the artist" example("Blue Town")
//	@Param 			sort query string false "criteria for sorting artist-searching results" example("-name", "name")
//	@Param 			page query int false "searching page" example(2)
//	@Param 			limit query int false "searching limit" example(10)
//	@Success		200	{object}	model.Artists
//	@Failure		400 "Bad Request"
//	@Failure		500 "Internal Server Error"
//	@Router			/artists [get]
func (ah *ArtistHandler) GetArtistWithFilter(c *gin.Context) {
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

	tracks, err := ah.repository.GetArtistsWithFilter(context.Background(), filter)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// DeleteArtist godoc
//
//		@Summary		Delete an artist
//		@Description	Delete an artist using ID
//		@Tags			artists
//		@Produce		json
//	 	@Param 			id path string true "Artist ID" example("3983a1d6-759b-4e5e-b307-7b7e06a05a85")
//		@Success		200	{object} response.DeleteArtistResponse
//		@Failure		400 "Bad Request"
//		@Failure		500 "Internal Server Error"
//		@Router			/artists/{id} [delete]
func (ah *ArtistHandler) DeleteArtist(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	err = ah.repository.DeleteArtist(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	delete_response := fmt.Sprintf("delete artist with id %v successfully", id)
	c.JSON(http.StatusOK, response.DeleteArtistResponse{Response: delete_response})
}

// UpdateArtistInformation godoc
//
//	@Summary		Update information of an artist
//	@Description	Update information of an artist
//	@Tags			artists
//	@Accept			json
//	@Produce		json
//	@Param 			artist body model.Artist true "artist information"
//	@Success		200	{object}	model.Artist
//	@Failure		400 "Bad Request"
//	@Failure		500 "Internal Server Error"
//	@Router			/artists [put]
func (ah *ArtistHandler) UpdateArtist(c *gin.Context) {
	artist := model.Artist{}
	if err := c.BindJSON(&artist); err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if err := ah.repository.UpdateArtist(context.Background(), &artist); err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusAccepted, artist)
}
