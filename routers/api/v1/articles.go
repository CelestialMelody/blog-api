package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetArticle
// @Title Get
// @Description get articles
// @Success 200 {object} models.Article
// @router / [get]
func GetArticle(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "get an article",
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
