package cache

import (
	"gin-gorm-practice/pkg/e"
	"strconv"
	"strings"
)

type Article struct {
	ID       int
	TagID    int
	State    int
	PageNum  int
	PageSize int
}

func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	key := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		key = append(key, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		key = append(key, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		key = append(key, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		key = append(key, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		key = append(key, strconv.Itoa(a.PageSize))
	}

	return strings.Join(key, "_")
}
