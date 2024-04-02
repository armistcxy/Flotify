package handler

import (
	"context"
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
	auth_manager helper.AuthManager
}

func NewUserHandler(repo repository.UserRepository) UserHandler {
	return UserHandler{
		repository:   repo,
		auth_manager: helper.NewAuthManager(),
	}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
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
	credential := helper.AuthCredential{
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
	// log.Println(token_string)
	// c.SetCookie("cookie", token_string, int(time.Now().UTC().Add(time.Minute*15).Unix()), "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"access token": access_token_string, "refresh token": refresh_token_string})
}

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
}
