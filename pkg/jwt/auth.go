package jwt

import (
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"blog-api/pkg/redis"
	"blog-api/pkg/util"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

const url = "http://localhost:8080/api/v1/login?token="

func backendLogin(c *gin.Context, token string) (*http.Request, error) {
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	// 后台登录更新token，本质上就是给login接口发送请求
	request, err := http.NewRequest("POST", url+token, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Logger.Error("login move forward error")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(map[string]interface{})
		appG := app.Gin{C: c}

		token := c.Query("token")
		if token == "" {
			log.Logger.Error(e.GetMsg(e.TokenEmpty))
			appG.Response(http.StatusUnauthorized, e.TokenEmpty, nil)
			// 终止后续操作
			c.Abort()
			return
		}

		// 判断 access token是否过期(2h)
		accTimeout, err0 := util.ValidToken(token)

		// token 未过期
		if !accTimeout {
			userID, err1 := util.GetUserIDFormToken(token)
			if err1 != nil {
				panic(err1)
			}
			c.Set("userID", userID)
			c.Next()
			return
		}

		// token 过期
		if accTimeout {
			// token过期或者解析token发生错误
			// 一般都是token过期
			log.Logger.Error("valid token err(accTimeout)", zap.Error(err0))
			log.Logger.Info("tern to valid refreshToken...")

			// refresh token 能否取出
			refreshToken, err1 := redis.RDB.Get(context.Background(), token).Result()

			if err1 != nil {
				log.Logger.Error("refresh token err(token wrong)", zap.Error(err1))
				log.Logger.Info("token is wrongful, re login please")
				data["token"] = token
				appG.Response(http.StatusUnauthorized, e.CheckTokenFail, data)
				c.Abort()
				return
			}

			// debug
			log.Logger.Debug("auth", zap.Any("token", token))
			log.Logger.Debug("auth", zap.Any("refresh token", refreshToken))

			// 可以取出30d token, 检查是否过期
			refreshTimeout, err2 := util.ValidToken(refreshToken)
			// 过期
			if refreshTimeout {
				//refreshToken出问题，表明用户三十天未登录，需要重新登录
				log.Logger.Error("valid refreshToken err:", zap.Error(err2))
				log.Logger.Info("user need login again")
				data["token"] = token
				appG.Response(http.StatusUnauthorized, e.CheckRefreshTokenTimeout, data)
				c.Abort()
				return
			}

			// refresh token 没过期
			//userId, err := util.GetUserIDFormToken(refreshToken)
			// 解析id出错 一般不能可能发生
			//if err != nil || userId == -1 {
			//	log.Logger.Error("parse token to get uid error:", zap.Error(err))
			//	log.Logger.Info("user need login again")
			//	appG.Response(http.StatusUnauthorized, e.CheckTokenFail, nil)
			//	c.Abort()
			//	return
			//}

			// refresh token 没过期
			userID, _ := util.GetUserIDFormToken(refreshToken)

			// 根据 refreshToken 更新 accessToken
			newToken, _ := util.GenerateToken(userID)

			// 一般不会出错
			//accessToken, err := util.GenerateToken(userID)
			//if err != nil {
			//	log.Logger.Error("create refresh token error:", zap.Error(err))
			//	log.Logger.Info("user need login again")
			//	appG.Response(http.StatusUnauthorized, e.CheckTokenFail, nil)
			//	c.Abort()
			//	return
			//}

			// 更新后，重新设置redis的key
			newRefreshToken, _ := util.GenerateRefreshToken(userID)

			// 一般不会出错
			//newRefreshToken, err := util.CreateRefreshToken(userId)
			//if err != nil {
			//	log.Logger.Error("creat ref token error:", zap.Error(err))
			//	log.Logger.Info("user need login again")
			//	appG.Response(http.StatusUnauthorized, e.CheckTokenFail, nil)
			//	c.Abort()
			//	return
			//}

			data["user_id"] = userID
			data["new_token"] = newToken
			data["new_refresh_token"] = newRefreshToken
			data["info"] = ""

			if err := redis.RDB.Set(context.Background(), token, newRefreshToken, 30*24*time.Hour).Err(); err != nil {
				log.Logger.Error("create redis refresh token error", zap.Error(err))
			} else {
				log.Logger.Info("redis set success")
			}

			// 更新成功
			//c.Set("userId", userID)
			//c.Next()
			//return

			//重新发起请求
			var request *http.Request
			var err3 error
			if request, err3 = backendLogin(c, newToken); err3 != nil {
				data["info"] = "reset request failed"
				log.Logger.Error("reset request failed", zap.Error(err3))
				appG.Response(http.StatusUnauthorized, e.ResetRequestFail, data)
				c.Abort()
				return
			}

			client := &http.Client{}
			post, err4 := client.Do(request)
			if post.StatusCode == 200 {
				//发送登录请求成功
				c.Set("userID", userID)
				c.Next()
				return
			} else {
				data["info"] = "backend login success"
				log.Logger.Error("restart request failed", zap.Error(err4))
				appG.Response(http.StatusUnauthorized, e.BackendLoginFail, data)
				c.Abort()
				return
			}
		}
	}
}
