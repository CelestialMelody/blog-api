package userSrv

import (
	"blog-api/internal/dao"
	model "blog-api/internal/models"
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"blog-api/pkg/redis"
	"blog-api/pkg/util"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	password string
}

type Req struct {
	Username string `form:"username" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

type Resp struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

// CheckExist
// 存在返回 true
func (u *User) CheckExist(username string) bool {
	// 检查是否存在
	_, err := dao.GetAuthorByUsername(username)
	if err != nil { // 未查询到
		return false
	}
	return true
}

func (u *User) GetAuthorInfo(username string) error {
	authorInTable, err := dao.GetAuthorByUsername(username)
	if err != nil {
		return errors.New(e.GetMsg(e.UserNotExist))
	}
	u.ID = authorInTable.ID
	u.Username = authorInTable.Username
	u.Email = authorInTable.Email
	u.password = authorInTable.Password
	return nil
}

func (u *User) Register(req Req) (Resp, error) {
	resp := Resp{}
	// 加密算法
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// 密码加密
	user := model.User{
		Username: req.Username,
		Password: string(hash),
	}

	userID, err := dao.Register(user)
	if err != nil {
		app.MarkError(err)
		return resp, err
	}

	// err == nil, id != -1
	token, _ := tokenSettings(userID, req.Username)
	resp.UserId = userID
	resp.Token = token
	return resp, nil
}

func tokenSettings(id int, username string) (string, error) {
	// acc token: 3h
	var token, refreshToken string
	var err error
	token, err = util.GenerateToken(id, username)
	if err != nil {
		return token, errors.New(e.GetMsg(e.GenerateTokenFail))
	}

	// ref token 30d
	refreshToken, err = util.GenerateRefreshToken(id, username)
	if err != nil {
		return token, errors.New(e.GetMsg(e.GenerateRefreshTokenFail))
	}

	// redis
	// key: 2h token
	// value 30d token
	// key live time: 30d
	if err := redis.RDB.Set(context.Background(),
		token, refreshToken,
		30*24*time.Hour).
		Err(); err != nil {
		log.Logger.Error("redis set error", zap.Error(err))
		return token, err
	} else {
		log.Logger.Debug("redis set success")
	}

	return token, nil
}

func (u *User) Login(req Req) (Resp, error) {
	resp := Resp{}
	if err := u.GetAuthorInfo(req.Username); err != nil {
		return resp, err
	}
	// 密码校验
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(req.Password))
	if err != nil {
		log.Logger.Error("password error", zap.Any("user", u))
		return resp, err
	}

	// err == nil, id != -1
	token, _ := tokenSettings(u.ID, req.Username)
	resp.UserId = u.ID
	resp.Token = token
	return resp, nil
}