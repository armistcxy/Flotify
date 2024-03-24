package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) (int, error) {
	var page int
	var err error
	page_string_form := c.Query("page")
	if page_string_form == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(page_string_form)
		if err != nil {
			return -1, err
		}
	}
	return page, nil
}

func GetLimit(c *gin.Context) (int, error) {
	var limit int
	var err error
	limit_string_form := c.Query("limit")
	if limit_string_form == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limit_string_form)
		if err != nil {
			return -1, err
		}
	}
	return limit, nil
}
