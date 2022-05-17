package api

import (
	"gin-gorm-practice/models/blogAuth"
	"gin-gorm-practice/pkg/app"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/util"
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
		appG.Response(http.StatusOK, e.INVALID_PARAMS, data)
		return
	}

	// 验证
	if err := blogAuth.CheckAuth(need.Username, need.Password); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, data)
		return
	}

	// 生成token
	if token, err := util.GenerateToken(need.Username, need.Password); err == nil {
		data["token"] = token
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		app.MarkError(err)
		appG.Response(http.StatusOK, e.ERROR_AUTH_TOKEN, data)
	}
}
