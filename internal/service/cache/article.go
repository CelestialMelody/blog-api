package cache

import (
	"blog-api/pkg/e"
	"strconv"
	"strings"
)

type Article struct {
	ID       int
	TagID    int
	PageNum  int
	PageSize int
}

func (a *Article) GetArticleKey() string {
	return e.CacheArticle + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	key := []string{
		e.CacheArticle,
		"LIST",
	}

	if a.ID > 0 {
		key = append(key, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		key = append(key, strconv.Itoa(a.TagID))
	}
	if a.PageNum > 0 {
		key = append(key, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		key = append(key, strconv.Itoa(a.PageSize))
	}

	return strings.Join(key, "_")
}
