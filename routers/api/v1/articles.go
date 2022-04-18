package v1

import (
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/models/blogArticle"
	"gin-gorm-practice/models/blogTag"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/util"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

// GetArticle
// @Summary Get a single article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code:200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	// 获取参数
	id := com.StrTo(c.Param("id")).MustInt()

	logger := zap.NewExample() // logger 和 valid 可以放在init函数中

	// 验证参数
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}

	if !valid.HasErrors() {
		if err := blogArticle.ExistArticleByID(id); err != nil {
			data = blogArticle.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			// 自己的日志
			logging.LoggoZap.Info(err.Key, zap.Any("message", err.Message))
			logger.Info(err.Key, zap.String("message", err.Message))
		}
	}

	c.JSON(200, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        data,
		"article_msg": "get an article",
	})
}

// GetArticles
// @Summary Get multiple articles
// @Produce  json
// @Param tag_id query int false "TagID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Param keyword query string false "Keyword"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	// 获取参数
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validator.New() // 使用 playground validator 包生成的验证器
	logger := zap.NewExample()

	// 验证参数
	type needValid struct {
		state int `validate:"oneof=0 1"`
		tagId int `validate:"min=1"`
	}

	need := &needValid{
		state: -1,
		tagId: -1,
	}

	if arg := c.Query("state"); arg != "" {
		need.state = com.StrTo(arg).MustInt()
		maps["state"] = need.state
	}

	if arg := c.Query("tag_id"); arg != "" {
		need.tagId = com.StrTo(arg).MustInt() // MustUint() ?
		maps["tag_id"] = need.tagId
	}

	code := e.INVALID_PARAMS

	// v0.4.1 validator貌似没有检验到必填项
	if err := valid.Struct(need); err == nil {
		code = e.SUCCESS
		data["lists"] = blogArticle.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = blogArticle.GetArticleTotalCount(maps)
	} else {
		code = e.ERROR_NOT_EXIST_ARTICLE
		logger.Info("validate error", zap.Any("error", err))
		// 自己的日志
		logging.LoggoZap.Info("validate error", zap.Any("message", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        data,
		"article_msg": "get articles",
	})
}

// AddArticle
// @Summary Add an article
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	// 获取参数
	type needValid struct {
		tagId     int    `validate:"min=1"`
		title     string `validate:"min=1,max=100"`
		desc      string `validate:"min=1,max=255"`
		content   string `validate:"min=1,max=65535"`
		createdBy string `validate:"min=1,max=100"`
		state     int    `validate:"oneof=0 1"`
	}

	need := needValid{
		tagId:     0, //
		title:     "",
		desc:      "",
		content:   "",
		createdBy: "",
		state:     -1,
	}

	need.tagId = com.StrTo(c.Query("tag_id")).MustInt() // 现在不考虑uint了// v0.3.1 uint -> 否则导致接口断言失败->AddArticle: TagID: data[ "tag_id"].(uint)
	need.title = c.Query("title")
	need.desc = c.Query("desc")
	need.content = c.Query("content")
	need.createdBy = c.Query("created_by")
	need.state = com.StrTo(c.Query("state")).MustInt()

	// 验证参数
	valid := validator.New() // 使用 playground validator 包生成的验证器
	logger := zap.NewExample()

	code := e.INVALID_PARAMS
	if err := valid.Struct(need); err == nil {
		if blogTag.ExistTagByID(need.tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = need.tagId
			data["title"] = need.title
			data["desc"] = need.desc
			data["content"] = need.content
			data["created_by"] = need.createdBy
			data["state"] = need.state
			if err := blogArticle.AddArticle(data); err == nil {
				code = e.SUCCESS
			} else {
				logger.Info("err: " + err.Error())
				// 自己的日志
				logging.LoggoZap.Info("validate error", zap.Any("message", err))
			}
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		logger.Info("err: " + err.Error())
		// 自己的日志
		logging.LoggoZap.Info("validate error", zap.Any("message", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "add an article",
	})
}

// EditArticle
// @Summary Edit an article
// @Produce  json
// @Param id body int true "ID"
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int true "State"
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{},"article_msg":"edit an article"}"
// @Router /api/v1/articles [put]
func EditArticle(c *gin.Context) {
	// 获取参数
	type needValid struct {
		id         int    `validate:"min=1"`
		tagID      int    `validate:"min=1"`
		title      string `validate:"min=1,max=100"`
		desc       string `validate:"min=1,max=255"`
		content    string `validate:"min=1,max=65535"`
		modifiedBy string `validate:"min=1,max=100"`
		state      int    `validate:"oneof=0 1"`
	}

	need := needValid{
		id:         -1,
		tagID:      -1,
		title:      "",
		desc:       "",
		content:    "",
		modifiedBy: "",
		state:      -1,
	}

	need.id = com.StrTo(c.Param("id")).MustInt()
	need.tagID = com.StrTo(c.Query("tag_id")).MustInt()
	need.title = c.Query("title")
	need.desc = c.Query("content")
	need.modifiedBy = c.Query("modified_by")
	need.state = com.StrTo(c.Query("state")).MustInt()

	logger := zap.NewExample()
	valid := validator.New()
	code := e.INVALID_PARAMS

	if err := valid.Struct(need); err == nil {
		if err := blogArticle.ExistArticleByID(need.id); err == nil {
			if blogTag.ExistTagByID(need.tagID) {
				data := make(map[string]interface{})
				data["tag_id"] = need.tagID
				data["title"] = need.title
				data["desc"] = need.desc
				data["content"] = need.content
				data["modified_by"] = need.modifiedBy
				data["state"] = need.state

				blogArticle.EditArticle(need.id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logger.Debug("err: " + err.Error())
			// 自己的日志
			logging.LoggoZap.Info("validate error", zap.Any("message", err))
		}
	} else {
		logger.Debug("err: " + err.Error())
		// 自己的日志
		logging.LoggoZap.Info("validate error", zap.Any("message", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "edit an article",
	})
}

// DeleteArticle
// @Summary Delete an article
// @Description Delete an article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	logger := zap.NewExample()
	valid := validator.New()

	code := e.INVALID_PARAMS

	if err := valid.Var(id, "min=1"); err == nil {
		if err := blogArticle.ExistArticleByID(id); err == nil {
			blogArticle.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logger.Info("err" + err.Error())
			// 自己的日志
			logging.LoggoZap.Info("validate error", zap.Any("message", err))
		}
	} else {
		logger.Info("err" + err.Error())
		// 自己的日志
		logging.LoggoZap.Info("validate error", zap.Any("message", err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "delete an article",
	})
}
