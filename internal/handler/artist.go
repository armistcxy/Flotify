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

type ArtistHandler struct {
	repository repository.ArtistRepository
}

func NewArtistHandler(repo repository.ArtistRepository) ArtistHandler {
	return ArtistHandler{
		repository: repo,
	}
}

// CreateArtist godoc
//	@Summary		Create an artist
//	@Description	Create a new artist
//	@Tags			artist
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Artist
//	@Failure		400
//	@Failure		500
//	@Router			/artist [post]
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

func (ah *ArtistHandler) GetArtistTracksByID(c *gin.Context) {
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

	tracks, err := ah.repository.GetTrackOfArtist(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	response := gin.H{
		"artist": artist,
		"tracks": tracks,
	}
	c.JSON(http.StatusOK, response)
}

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

	c.JSON(http.StatusOK, gin.H{"message": "delete successfully"})
}

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
