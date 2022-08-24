package controller

import (
	"blog-api/internal/dao"
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 用户注册
// @Summary 用户注册
// @Produce  json
// @Param name query string true "Name"
// @Param password query string true "Password"
// @Param email query string true "Email"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/register [post]
func Register(c *gin.Context) {
	type needValid struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	var need needValid
	appG := app.Gin{C: c}

	need.Username = c.Query("username")
	need.Password = c.Query("password")

	// 验证并处理参数
	if err := valid.Struct(need); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	if err := dao.CheckAuth(need.Username, need.Password); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.ErrorNotExistUser, nil)
		// 创建用户
		err := dao.Register(need.Username, need.Password)
		if err != nil {
			app.MarkError(err)
			appG.Response(http.StatusOK, e.ErrorRegisterFail, nil)
			return
		}
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
