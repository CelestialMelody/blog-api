package tagSrv

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

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string

	PageNum  int
	PageSize int
}

type AddReq struct {
	Name      string `binding:"required,max=100"`
	CreatedBy string `binding:"required,max=100"`
}

type EditReq struct {
	ID         int    `binding:"required,min=1"`
	Name       string `binding:"required,max=100"`
	ModifiedBy string `bingding:"required,max=100"`
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (t *Tag) ExistByName() error {
	return dao.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() error {
	return dao.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return dao.AddTag(t.Name, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	return dao.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return dao.DeleteTag(t.ID)
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}

	return maps
}

func (t *Tag) Count() (int64, error) {
	return dao.GetTagCount(t.getMaps())
}

func (t *Tag) GetAll() ([]model.Tag, error) {
	ctx := context.Background()
	tagCache := cache.Tag{
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := tagCache.GetKey()
	val, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		log.Logger.Error("get tag from redis error", zap.Error(err))
	} else {
		var tags []model.Tag
		_ = json.Unmarshal([]byte(val), &tags)
		return tags, nil
	}

	tags, err := dao.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	tagsBytes, _ := json.Marshal(tags)

	if err := redis.RDB.Set(ctx, key, tagsBytes, time.Hour).Err(); err != nil {
		log.Logger.Error("set tag to redis error", zap.Error(err))
	}

	return tags, nil
}
