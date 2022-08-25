package controller

import (
	userSrv "blog-api/internal/service/user"
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

// Register 用户注册
// @Summary 用户注册
// @Produce  json
// @Param username query string true "Name"
// @Param password query string true "Password"
// @Param email query string true "Email"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /register [post]
func Register(c *gin.Context) {
	appG := app.Gin{C: c}
	var data = make(map[string]interface{})
	var req userSrv.Req
	var resp userSrv.Resp
	var err error
	var u = userSrv.User{}

	// 验证并处理参数
	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	// 查询是否存在
	if ok := u.CheckExist(req.Username); ok {
		appG.Response(http.StatusOK, e.UsernameExist, data)
		return
	}

	// 注册
	resp, err = u.Register(req)
	if err != nil {
		log.Logger.Error("register fail", zap.Error(err))
		appG.Response(http.StatusOK, e.RegisterFail, nil)
		return
	}

	data["uid"] = resp.UserId
	data["token"] = resp.Token
	appG.Response(http.StatusOK, e.Success, data)
}

// GetAuthorInfo 获取作者信息
// @Summary 获取作者信息
// @Description 获取作者信息
// @Tags Auth
// @Accept json
// @Produce json
// @Param authorID query string true "author_id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /authorInfo [get]
func GetAuthorInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	data := make(map[string]interface{})
	var u = userSrv.User{}

	username := c.Query("username")
	token := c.Query("token")

	// 验证
	if err := u.GetAuthorInfo(username); err != nil {
		log.Logger.Error("get author info fail", zap.Error(err))
		data["author_id"] = u.ID
		data["username"] = u.Username
		data["email"] = u.Email
		data["token"] = token
		appG.Response(http.StatusOK, e.CheckTokenFail, data)
		return
	}
}

// Login 用户登录
// @Summary 用户登录
// @Produce  json
// @Param username query string true "Name"
// @Param password query string true "Password"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /login [post]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	data := make(map[string]interface{})
	var req = userSrv.Req{}
	var resp = userSrv.Resp{}
	var u = userSrv.User{}
	var err error

	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	// 查询是否存在
	if ok := u.CheckExist(req.Username); ok {
		appG.Response(http.StatusOK, e.UsernameExist, data)
		return
	}

	resp, err = u.Login(req)
	if err != nil {
		log.Logger.Error("login fail", zap.Error(err))
		appG.Response(http.StatusOK, e.LoginFail, nil)
		return
	}

	data["uid"] = resp.UserId
	data["token"] = resp.Token
	appG.Response(http.StatusOK, e.Success, data)
}
