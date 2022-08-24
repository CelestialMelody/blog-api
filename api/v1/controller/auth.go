package controller

import (
	"blog-api/internal/dao"
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	valid = validator.New()
)

// GetAuth 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags Auth
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /blogAuth [get]
func GetAuth(c *gin.Context) {
	type needValid struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// 获取参数
	need := needValid{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}
	appG := app.Gin{C: c}
	data := make(map[string]interface{})

	// 验证并处理参数
	if err := valid.Struct(need); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.InvalidParams, data)
		return
	}

	// 验证
	if err := dao.CheckAuth(need.Username, need.Password); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.ErrorAuthCheckTokenFail, data)
		return
	}

	// 生成token
	if token, err := util.GenerateToken(need.Username, need.Password); err == nil {
		data["token"] = token
		appG.Response(http.StatusOK, e.Success, data)
	} else {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.ErrorAuthToken, data)
	}
}
