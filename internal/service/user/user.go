package userSrv

import (
	"blog-api/internal/dao"
	model "blog-api/internal/models"
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
	UserID int    `json:"user_id"`
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
		log.Logger.Error("register error", zap.Error(err))
		return resp, err
	}

	// err == nil, id != -1
	token, _ := tokenSettings(userID)
	resp.UserID = userID
	resp.Token = token
	return resp, nil
}

func tokenSettings(id int) (string, error) {
	// acc token: 3h
	var token, refreshToken string
	token, _ = util.GenerateToken(id)

	//token, err = util.GenerateToken(id)
	// 一般不会失败
	//if err != nil {
	//	return token, errors.New(e.GetMsg(e.GenerateTokenFail))
	//}

	// ref token 30d
	refreshToken, _ = util.GenerateRefreshToken(id)

	// debug
	log.Logger.Debug("login/register", zap.Any("token", token))
	log.Logger.Debug("login/register", zap.Any("refresh token", refreshToken))

	//refreshToken, err = util.GenerateRefreshToken(id)
	// 一般不会失败
	//if err != nil {
	//	return token, errors.New(e.GetMsg(e.GenerateTokenFail))
	//}

	// redis
	// key: 2h token
	// value 30d token
	// key live time: 30d
	if err := redis.RDB.Set(context.Background(),
		token, refreshToken,
		30*24*time.Hour).
		Err(); err != nil {
		log.Logger.Error("login/ redis set error", zap.Error(err))
		return token, err
	} else {
		log.Logger.Info("login/ redis set success")
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
	token, _ := tokenSettings(u.ID)
	resp.UserID = u.ID
	resp.Token = token
	return resp, nil
}
