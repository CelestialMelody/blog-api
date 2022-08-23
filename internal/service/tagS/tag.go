package tagS

import (
	"blog-api/internal/dao"
	model "blog-api/internal/models"
	"blog-api/internal/service/cache"
	"blog-api/pkg/app"
	"blog-api/pkg/redis"
	jsoniter "github.com/json-iterator/go"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (t *Tag) ExistByName() error {
	return dao.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() error {
	return dao.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return dao.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

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

	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) Count() (int64, error) {
	return dao.GetTagCount(t.getMaps())
}

func (t *Tag) GetAll() ([]model.Tag, error) {
	tagCache := cache.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := tagCache.GetKey()
	if redis.Exists(key) {
		data, err := redis.Get(key)
		if err != nil {
			app.MarkError(err)
		} else {
			var tags []model.Tag
			_ = json.Unmarshal(data, &tags)
			return tags, nil
		}
	}

	tags, err := dao.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	if err := redis.Set(key, tags, 3600); err != nil {
		app.MarkError(err)
	}

	return tags, nil
}
