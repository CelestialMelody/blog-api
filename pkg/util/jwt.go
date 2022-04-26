package util

import (
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

var jwtSecret []byte
var logger = zap.NewExample().Sugar()

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//func init() {
//jwtSecret = []byte(setting.GlobalConfig.GetString("app.jwt_secret"))
// 方式2
//sec, err := setting.Cfg.GetSection("app")
//if err != nil {
//	logger.Error("init jwt secret error", zap.Error(err))
//	return
//}
//	jwtSecret = []byte(sec.Key("jwt_secret").String())
//	logger.Info("jetSecret", zap.String("jetSecret", string(jwtSecret)))
//}

// GenerateToken 生成token
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
	//token, err := tokenClaims.SignedString(jwtSecret)
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
	//jwtSecret = []byte(setting.GlobalConfig.GetString("app.jwt_secret"))
	token, err := tokenClaims.SignedString(jwtSecret)
	logging.LoggoZap.Info("token", zap.String("token", token))
	if err != nil {
		logging.LoggoZap.Error("generate token error", zap.Error(err))
	}

	// 报错了
	//claims := Claims{
	//	Username: username,
	//	Password: password,
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: expireTime.Unix(),
	//		Issuer:    "gin-blog",
	//	},
	//}

	return token, err
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	// 解析token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	// debug
	logger.Debug("token", zap.String("token", token), zap.Any("err", err))

	if tokenClaims != nil {
		// 获取自定义的claims
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // 校验token
			return claims, nil
		}
	}
	return nil, err
}
