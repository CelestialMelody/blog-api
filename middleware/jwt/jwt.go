package jwt

import (
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		code := e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			// 验证token
			claims, err := util.ParseToken(token)
			if err != nil { // token 校验失败
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt { // token 过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
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

			// 终止后续操作
			c.Abort()
			return
		}
		// 执行后续操作
		c.Next()
	}
}
