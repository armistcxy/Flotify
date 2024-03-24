package repository

import (
	"fmt"
	"strings"
)

type Filter struct {
	Props  map[string]any
	Page   int
	Limit  int
	SortBy []string
}

func (f Filter) GetOffSet() int {
	return f.Limit * (f.Page - 1)
}

func (f Filter) GetSortCriteria() string {
	sort_criteria := ""
	for _, criteria := range f.SortBy {
		field := strings.TrimPrefix(criteria, "-")
		dir := "ASC"
		if strings.HasPrefix(criteria, "-") {
			dir = "DESC"
		}
		sort_criteria = fmt.Sprintf("%s %s %s,", sort_criteria, field, dir)
	}
	return sort_criteria
}
