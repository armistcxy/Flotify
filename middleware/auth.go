package middleware

import (
	"errors"
	"flotify/internal/auth"
	"flotify/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

func AuthRequest(auth_manager auth.AuthManager) gin.HandlerFunc {

	return func(c *gin.Context) {
		token_string := c.Request.Header.Get("Authorization")
		if token_string == "" {
			err := errors.New("token nonexist")
			helper.ErrorResponse(c, err, http.StatusUnauthorized)
		}
		token := token_string[len("Bearer "):]

		id_string_form := c.Params.ByName("id")
		id, err := uuid.FromString(id_string_form)
		if err != nil {
			helper.ErrorResponse(c, err, http.StatusBadRequest)
			return
		}

		if err := auth_manager.VerifyJWT(token, id); err != nil {
			helper.ErrorResponse(c, err, http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
