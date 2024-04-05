package handler

import (
	"context"
	"flotify/internal/auth"
	"flotify/internal/custom_error"
	"flotify/internal/helper"
	"flotify/internal/model"
	"flotify/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type UserHandler struct {
	repository   repository.UserRepository
	auth_manager auth.AuthManager
}

func NewUserHandler(repo repository.UserRepository, auth_manager auth.AuthManager) UserHandler {
	return UserHandler{
		repository:   repo,
		auth_manager: auth_manager,
	}
}

// CreateUser godoc
// @Summary RegisterUser
// @Accept json
// @Produce json
// @Description Register user
// @Success 200 {response} model.user
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Tag users
// @Router /users/register [POST]
func (uh *UserHandler) RegisterUser(c *gin.Context) {
	type RequestUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	request_user := RequestUser{}
	err := c.BindJSON(&request_user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if len(request_user.Password) < 8 || len(request_user.Password) > 64 {
		helper.ErrorResponse(c, &custom_error.PasswordLengthError{}, http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: request_user.Username,
		Email:    request_user.Email,
		Password: request_user.Password,
	}

	user, err = uh.repository.CreateUser(context.Background(), user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

// LoginUser godoc
// @Summary Login user
// @Description Login user with credentials (email and password)
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/login [POST]
func (uh *UserHandler) LoginUser(c *gin.Context) {
	type RequestUser struct {
		Email    string
		Password string
	}

	request_user := RequestUser{}
	err := c.BindJSON(&request_user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if len(request_user.Password) < 8 || len(request_user.Password) > 64 {
		helper.ErrorResponse(c, &custom_error.PasswordLengthError{}, http.StatusBadRequest)
		return
	}
	var id *uuid.UUID
	if id, err = uh.repository.UserLogin(context.Background(), request_user.Email, request_user.Password); err != nil {
		switch err := err.(type) {
		case custom_error.MismatchError:
			helper.ErrorResponse(c, err, http.StatusBadRequest)
			return
		default:
			helper.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
	}

	// create token
	credential := auth.AuthCredential{
		ID: *id,
	}
	access_token_string, err := uh.auth_manager.GenerateJWT(&credential, time.Minute*15)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
	}
	refresh_token_string, err := uh.auth_manager.GenerateJWT(&credential, time.Hour*24*7)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}
	err = uh.auth_manager.StoreRefreshToken(*id, refresh_token_string)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
	}
	// log.Println(token_string)
	// c.SetCookie("cookie", token_string, int(time.Now().UTC().Add(time.Minute*15).Unix()), "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"access token": access_token_string, "refresh token": refresh_token_string})
}

// ViewUserInformation godoc
// @Summary View user information
// @Description View user information with user ID
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} model.user
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/{id} [GET]
func (uh *UserHandler) ViewInformation(c *gin.Context) {
	id_string_form := c.Params.ByName("id")

	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	user, err := uh.repository.GetUserByID(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

// ModifyUserInformation godoc
// @Summary Modify user information
// @Description Modify user information with user ID
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} model.user
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/{id} [PUT]
func (uh *UserHandler) ModifyInformation(c *gin.Context) {
	id_string_form := c.Params.ByName("id")
	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	type RequestUser struct {
		Username string `json:"username"`
	}

	request_user := RequestUser{}
	err = c.BindJSON(&request_user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	user := model.User{
		ID:       id,
		Username: request_user.Username,
	}
	err = uh.repository.UpdateUserInfo(context.Background(), &user)
	if err != nil {
		switch err := err.(type) {
		case custom_error.DuplicateUsernameError:
			helper.ErrorResponse(c, err, http.StatusBadRequest)
			return
		default:
			helper.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusAccepted, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user with user ID
// @Param id string path true
// @Success 200
// @Failure 400 "Bad request"
// @Failure 401 "Authorization required"
// @Failure 500 "Internal server error"
// @Router /users/{id} [DELETE]
func (ur *UserHandler) DeleteUser(c *gin.Context) {
	// require authentication and authorization

	id_string_form := c.Params.ByName("id")
	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	err = ur.repository.DeleteUser(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete successfully"})
}

func (ur *UserHandler) GetFollowArtist(c *gin.Context) {
	id_string_form := c.Params.ByName("id")
	id, err := uuid.FromString(id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	artists, err := ur.repository.GetFollowArtist(context.Background(), id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, artists)
}

func (ur *UserHandler) FollowArtist(c *gin.Context) {
	user_id_string_form := c.Params.ByName("id")
	user_id, err := uuid.FromString(user_id_string_form)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	type RequestArtist struct {
		ID string `json:"id"`
	}
	request_artist := RequestArtist{}
	err = c.BindJSON(&request_artist)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}
	artist_id, err := uuid.FromString(request_artist.ID)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}
	err = ur.repository.FollowArtist(context.Background(), user_id, artist_id)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "follow artist successfully"})
}


