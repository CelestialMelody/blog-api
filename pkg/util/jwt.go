package util

import (
	"blog-api/conf"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

var jwtSecret []byte

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func createToken(id int, t time.Duration) (string, error) {
	var token string
	var err error
	// 当前时间
	nowTime := time.Now()
	// 设置过期时间
	expireTime := nowTime.Add(t)
	// 声明
	claims := Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "blog-api",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	jwtSecret = []byte(conf.AppConfig.JwtSecret)
	token, err = tokenClaims.SignedString(jwtSecret)

	log.Logger.Debug("token", zap.String("token", token))

	if err != nil {
		log.Logger.Error(e.GetMsg(e.GenerateTokenFail))
	}

	return token, err
}

func GenerateToken(id int) (string, error) {
	//t := 3 * time.Hour

	// debug

	t := 60 * time.Second

	return createToken(id, t)
}

func GenerateRefreshToken(id int) (string, error) {
	t := 30 * 24 * time.Hour
	return createToken(id, t)
}

func ParseToken(token string) (*Claims, error) {
	// 解析token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if tokenClaims != nil {
		// 获取自定义的claims
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // 校验token
			return claims, nil
		}
	}
	return nil, err
}

// GetUserIDFormToken ParseToken 验证用户token
// id int: 用户id 如果没有解析出，默认为-1
// err error: 错误
func GetUserIDFormToken(token string) (int, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return -1, err
	}
	return claims.ID, nil
}

// ValidToken 校验token是否过期
// bool: 是否过期 true 过期
// error: 解析是否成功 nil
func ValidToken(token string) (bool, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return true, err
	}
	expiresTime := claims.ExpiresAt
	now := time.Now().Unix()
	if now > expiresTime {
		//token 过期了
		return true, nil
	}
	return false, nil
}
