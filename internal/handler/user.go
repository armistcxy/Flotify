package handler

import (
	"context"
	"flotify/internal/helper"
	"flotify/internal/model"
	"flotify/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) UserHandler {
	return UserHandler{
		repository: repo,
	}
}

func (ur *UserHandler) CreateUser(c *gin.Context) {
	type RequestUser struct {
		Name     string
		Password string
	}

	request_user := RequestUser{}
	err := c.BindJSON(&request_user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: request_user.Name,
		Password: request_user.Password,
	}

	user, err = ur.repository.CreateUser(context.Background(), user)
	if err != nil {
		helper.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}
