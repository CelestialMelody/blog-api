package jwt

import (
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var logger = zap.NewExample().Sugar()

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
				//logging.Error("ParseTokenFailed", zap.Error(err)) // demo 测试自己的日志输出
				logger.Error("token 校验失败", zap.String("err", err.Error()))
				logging.LoggoZap.Error("token 校验失败", zap.String("err", err.Error()))
			} else if time.Now().Unix() > claims.ExpiresAt { // token 过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				//logging.Error("TokenOutOfDate", zap.Error(err)) // demo 测试自己的日志输出
				logger.Error("token 过期", zap.String("err", err.Error()))
				logging.LoggoZap.Error("token 过期", zap.String("err", err.Error()))
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

			//logging.Error("JWT", zap.String("err", e.GetMsg(code)))          // demo 测试自己的日志输出
			logging.LoggoZap.Error("JWT", zap.String("err", e.GetMsg(code))) // demo 测试自己的日志输出
			logger.Error("token 校验失败", zap.String("err", e.GetMsg(code)))

			// 终止后续操作
			c.Abort()
			return
		}
		// 执行后续操作
		c.Next()
	}
}
