package util

import (
	"gin-gorm-practice/conf"
	"gin-gorm-practice/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	// 当前时间
	nowTime := time.Now()
	// 设置过期时间
	expireTime := nowTime.Add(3 * time.Hour) // 3小时过期
	// 声明
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-gorm-practice",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	jwtSecret = []byte(conf.AppSetting.JwtSecret)
	token, err := tokenClaims.SignedString(jwtSecret)
	log.Logger.Info("token", zap.String("token", token))
	if err != nil {
		log.Logger.Error("generate token error", zap.Error(err))
	}

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	// 解析token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	log.Logger.Debug("token", zap.String("token", token), zap.Any("err", err))

	if tokenClaims != nil {
		// 获取自定义的claims
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // 校验token
			return claims, nil
		}
	}
	return nil, err
}
