package jwt

import (
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"blog-api/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		code := e.Success

		token := c.Query("token")
		if token == "" {
			code = e.InvalidParams
		} else {
			//claims, err := util.ParseToken(token)
			_, err := util.ParseToken(token)
			if err != nil { // token 校验失败
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired: // time.Now().Unix() > claims.ExpiresAt
					code = e.UserCheckTokenTimeout
					app.MarkError(err)
				default:
					code = e.UserCheckTokenFail
					app.MarkError(err)
				}
			}
		}

		// 失败返回结果
		if code != e.Success {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  code,
				"msg":   e.GetMsg(code),
				"data":  data,
				"token": token,
			})
			log.Logger.Error(e.GetMsg(code))
			// 终止后续操作
			c.Abort()
			return
		}

		// 执行后续操作
		c.Next()
	}
}
