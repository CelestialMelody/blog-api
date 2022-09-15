package cache

import (
	"blog-api/pkg/e"
	"strconv"
	"strings"
)

type Tag struct {
	ID   int
	Name string

	PageNum  int
	PageSize int
}

func (t *Tag) GetKey() string {
	key := []string{
		e.CacheTag,
		"LIST",
	}

	// t.ID ?

	if t.Name != "" {
		key = append(key, t.Name)
	}
	if t.PageNum > 0 {
		key = append(key, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		key = append(key, strconv.Itoa(t.PageSize))
	}

	return strings.Join(key, "_")
}
