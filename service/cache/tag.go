package cache

import (
	"gin-gorm-practice/pkg/e"
	"strconv"
	"strings"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

func (t *Tag) GetKey() string {
	key := []string{
		e.CACHE_TAG,
		"LIST",
	}

	// t.ID ?

	if t.Name != "" {
		key = append(key, t.Name)
	}
	if t.State >= 0 {
		key = append(key, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		key = append(key, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		key = append(key, strconv.Itoa(t.PageSize))
	}

	return strings.Join(key, "_")
}
