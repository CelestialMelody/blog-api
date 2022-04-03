package v1

import (
	"gin-gorm-practice/models/blogArticle"
	"gin-gorm-practice/pkg/e"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
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

	logger := zap.NewExample()

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
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// GetArticles
// @Title Get
// @Description get articles
// @Success 200 {object} models.Article
// @router / [get]
func GetArticles(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "get articles",
	})
}

// AddArticle
// @Title Post
// @Description post an article
// @Success 200 {object} models.Article
// @router / [post]
func AddArticle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "add an article",
	})
}

// EditArticle
// @Title Put
// @Description edit an article
// @Success 200 {object} models.Article
// @router / [put]
func EditArticle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "edit an article",
	})
}

// DeleteArticle
// @Title Delete
// @Description delete an article
// @Success 200 {object} models.Article
// @router / [delete]
func DeleteArticle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "delete an article",
	})
}
