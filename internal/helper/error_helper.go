package helper

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, err error, code int) {
	c.Error(err)
	c.AbortWithStatusJSON(code, gin.H{"status": false, "message": err.Error()})
}
