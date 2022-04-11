package v1

import (
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/models/blogArticle"
	"gin-gorm-practice/models/blogTag"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/util"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

// GetArticle
// @Title Get
// @Description get articles
// @Success 200 {object} models.Article
// @router / [get]
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
			logger.Info(err.Key, zap.String("message", err.Message))
			//logger.Error(err.Key, zap.String("message", err.Message))
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
// @Title Get
// @Description get articles
// @Success 200 {object} models.Article
// @router / [get]
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
	var need needValid
	need.state = -1
	need.tagId = -1

	if arg := c.Query("state"); arg != "" {
		need.state = com.StrTo(arg).MustInt()
		maps["state"] = need.state
	}

	if arg := c.Query("tag_id"); arg != "" {
		need.tagId = com.StrTo(arg).MustInt() // MustUint() ?
		maps["tag_id"] = need.tagId
	}

	code := e.INVALID_PARAMS

	if err := valid.Struct(need); err == nil {
		code = e.SUCCESS
		data["lists"] = blogArticle.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = blogArticle.GetArticleTotalCount(maps)
	} else {
		logger.Info(err.Error())
		code = e.ERROR_NOT_EXIST_ARTICLE
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        data,
		"article_msg": "get articles",
	})
	// 没看到返回结果 只能手动加了
	logger.Info("GetArticles", zap.Any("data", data))
}

// AddArticle
// @Title Post
// @Description post an article
// @Success 200 {object} models.Article
// @router / [post]
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
			}
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		logger.Info("err: " + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "add an article",
	})
}

// EditArticle
// @Title Put
// @Description edit an article
// @Success 200 {object} models.Article
// @router / [put]
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
		}
	} else {
		logger.Debug("err: " + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "edit an article",
	})
}

// DeleteArticle
// @Title Delete
// @Description delete an article
// @Success 200 {object} models.Article
// @router / [delete]
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
		}
	} else {
		logger.Info("err" + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         e.GetMsg(code),
		"data":        make(map[string]interface{}),
		"article_msg": "delete an article",
	})
}
