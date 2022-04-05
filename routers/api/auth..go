package api

import (
	"gin-gorm-practice/models/blogAuth"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

var (
	logger = zap.NewExample()
	valid  = validator.New()
)

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

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	// 验证并处理参数
	if err := valid.Struct(need); err == nil {
		// 验证
		if isExist := blogAuth.CheckAuth(need.Username, need.Password); isExist == true {
			// 生成token
			if token, errAuth := util.GenerateToken(need.Username, need.Password); errAuth == nil {
				data["token"] = token
				code = e.SUCCESS
			} else {
				code = e.ERROR_AUTH_TOKEN
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		logger.Info("validate error", zap.Error(err))
	}

	// 查看
	logger.Info("auth", zap.Any("data", data), zap.Any("code", code))

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     e.GetMsg(code),
		"data":    data,
		"message": "pong",
	})
}
