package jwt

import (
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/log"
	"gin-gorm-practice/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		code := e.SUCCESS

		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			//claims, err := util.ParseToken(token)
			_, err := util.ParseToken(token)
			if err != nil { // token 校验失败
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired: // time.Now().Unix() > claims.ExpiresAt
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					log.Logger.Error("token 过期", zap.String("err", err.Error()))
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
					log.Logger.Error("token 校验失败", zap.String("err", err.Error()))
				}
			}
		}

		// 失败返回结果
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  code,
				"msg":   e.GetMsg(code),
				"data":  data,
				"token": token,
			})
			log.Logger.Error("JWT", zap.String("err", e.GetMsg(code)))
			// 终止后续操作
			c.Abort()
			return
		}

		// 执行后续操作
		c.Next()
	}
}
