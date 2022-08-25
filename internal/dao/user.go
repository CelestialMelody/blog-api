package dao

import (
	model "blog-api/internal/models"
	"blog-api/pkg/log"
	"blog-api/pkg/mysql"
	"go.uber.org/zap"
)

// GetAuthorByUsername
// 通过用户名查询是否存在;
// 不存在返回 user(不为nil但值为0), error;
// 存在返回 user, nil
func GetAuthorByUsername(username string) (*model.User, error) {
	var u *model.User
	// 使用first查询 没有找到 err
	// 若使用find查询 err 为 nil
	// 实际上无论是find 还是 first查询后 u都不为nil
	err := mysql.DB.Debug().
		Where("username = ?", username).
		First(&u).Error

	// debug
	log.Logger.Debug("get user by username", zap.Any("user", u))

	return u, err
}

func GetUsernameByID(id int) (string, error) {
	var u *model.User
	err := mysql.DB.Debug().
		Where("id = ?", id).
		First(&u).Error

	return u.Username, err
}

func Register(u model.User) (int, error) {
	err := mysql.DB.Create(&u).Error
	if err != nil {
		return -1, err
	}
	return u.ID, nil
}
