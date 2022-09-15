package articleSrv

import (
	"blog-api/internal/dao"
	model "blog-api/internal/models"
	"blog-api/internal/service/cache"
	"blog-api/pkg/log"
	"blog-api/pkg/redis"
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

type EditReq struct {
	ID            int    `validate:"min=1"`
	TagID         int    `validate:"min=1"`
	Title         string `validate:"min=1,max=100"`
	Desc          string `validate:"min=1,max=255"`
	Content       string `validate:"min=1,max=65535"`
	ModifiedBy    string `validate:"min=1,max=100"`
	CoverImageUrl string
}

type AddReq struct {
	TagID         int    `binding:"min=1"`
	Title         string `binding:"min=1,max=100"`
	Desc          string `binding:"min=1,max=255"`
	Content       string `binding:"min=1,max=65535"`
	CreatedBy     string `binding:"min=1,max=100"`
	CoverImageUrl string
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
	}

	err := dao.AddArticle(article)
	return err

}

func (a *Article) Edit() error {
	return dao.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"modified_by":     a.ModifiedBy,
	})
}

// 先从缓存中取; 仅对redis的error进行处理
// 缓存中没有再从数据库中取

func (a *Article) Get() (*model.Article, error) {
	articleCache := cache.Article{ID: a.ID}
	key := articleCache.GetArticleKey()
	ctx := context.Background()
	val, err := redis.RDB.Get(ctx, key).Result()

	if err != nil {
		log.Logger.Error("get article from redis error", zap.Error(err))
	} else {
		var article *model.Article
		_ = json.Unmarshal([]byte(val), &article)
		return article, nil
	}

	article, err := dao.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	articleBytes, err := json.Marshal(article)

	if err := redis.RDB.Set(ctx, key, articleBytes, time.Hour).Err(); err != nil {
		log.Logger.Error("set article to redis error", zap.Error(err))
	}
	return article, nil
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}

func (a *Article) GetAll() ([]*model.Article, error) {
	ctx := context.Background()
	articlesCache := cache.Article{
		TagID:    a.TagID,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	key := articlesCache.GetArticlesKey()
	val, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		log.Logger.Error("get articles from redis error", zap.Error(err))
	} else {
		var articles []*model.Article
		_ = json.Unmarshal([]byte(val), &articles)
		return articles, nil
	}

	articles, err := dao.GetArticleLists(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	articlesBytes, err := json.Marshal(articles)

	if err := redis.RDB.Set(ctx, key, articlesBytes, time.Hour).Err(); err != nil {
		log.Logger.Error("set articles to redis error", zap.Error(err))
	}
	return articles, nil
}

func (a *Article) Delete() error {
	return dao.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() error {
	return dao.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int64, error) {
	return dao.GetArticleTotalCount(a.getMaps())
}
