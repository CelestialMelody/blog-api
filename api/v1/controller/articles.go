package controller

import (
	"blog-api/conf"
	"blog-api/internal/dao"
	articleSrv "blog-api/internal/service/article"
	tagSrv "blog-api/internal/service/tag"
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"blog-api/pkg/util"
	"blog-api/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

// GetArticle
// @Summary Get a single blogArticle
// @Description 获取文章
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{C: c}

	if err := validate.Var(id, "required,min=1"); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	if err := dao.ExistArticleByID(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.NotExistArticle, nil)
	}

	articleService := articleSrv.Article{ID: id}
	article, err := articleService.Get()
	if err != nil {
		log.Logger.Error("get article fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.GetArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, article)
}

// GetArticleLists
// @Summary Get multiple articles
// @Description 获取多篇文章
// @Tags 文章
// @Produce json
// @Param tag_id query int false "标签ID"
// @Param state query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticleLists(c *gin.Context) {
	// 验证参数
	type needValid struct {
		state int `validate:"oneof=0 1"`
		tagId int `validate:"min=1"`
	}
	var need needValid
	appG := app.Gin{C: c}

	need.state = com.StrTo(c.Query("state")).MustInt()
	need.tagId = com.StrTo(c.Query("tag_id")).MustInt()

	if err := validate.Struct(need); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := articleSrv.Article{
		TagID:    need.tagId,
		State:    need.state,
		PageNum:  util.GetPage(c),
		PageSize: conf.AppConfig.PageSize,
	}

	articleLists, err := articleService.GetAll()
	if err != nil {
		log.Logger.Error("get article lists fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.GetArticleListFail, nil)
		return
	}

	total, err := articleService.Count()
	if err != nil {
		log.Logger.Error("get article lists count fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.GetArticleCountFail, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articleLists
	data["total"] = total

	appG.Response(http.StatusOK, e.Success, data)
}

// AddArticle
// @Summary Add a blogArticle
// @Description 添加文章
// @Tags 文章
// @Produce json
// @Param tag_id query int true "标签ID"
// @Param title query string true "标题"
// @Param desc query string true "描述"
// @Param content query string true "内容"
// @Param created_by query string true "创建人"
// @Param state query int true "状态"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	// 获取参数
	type needValid struct {
		tagId         int    `validate:"min=1"`
		title         string `validate:"min=1,max=100"`
		desc          string `validate:"min=1,max=255"`
		content       string `validate:"min=1,max=65535"`
		createdBy     string `validate:"min=1,max=100"`
		state         int    `validate:"oneof=0 1"`
		coverImageUrl string `validate:"min=1,max=255"`
	}

	var need needValid
	appG := app.Gin{C: c}

	need.tagId = com.StrTo(c.Query("tag_id")).MustInt()
	need.title = c.Query("title")
	need.desc = c.Query("desc")
	need.content = c.Query("content")
	need.createdBy = c.Query("created_by")
	need.state = com.StrTo(c.Query("state")).MustInt()
	need.coverImageUrl = c.Query("cover_image_url")

	if err := validate.Struct(need); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagService := tagSrv.Tag{ID: need.tagId}
	if err := tagService.ExistByID(); err != nil {
		log.Logger.Error("tag not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistTag, nil)
		return
	}

	articleService := articleSrv.Article{
		TagID:         need.tagId,
		Title:         need.title,
		Desc:          need.desc,
		Content:       need.content,
		CoverImageUrl: need.coverImageUrl,
		CreatedBy:     need.createdBy,
		State:         need.state,
	}

	if err := articleService.Add(); err != nil {
		log.Logger.Error("add article fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.AddArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// EditArticle
// @Summary Update a blogArticle
// @Description 更新文章
// @Tags 文章
// @Produce json
// @Param id path int true "ID"
// @Param tag_id query int true "标签ID"
// @Param title query string true "标题"
// @Param desc query string true "描述"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	// 获取参数
	type needValid struct {
		id            int    `validate:"min=1"`
		tagID         int    `validate:"min=1"`
		title         string `validate:"min=1,max=100"`
		desc          string `validate:"min=1,max=255"`
		content       string `validate:"min=1,max=65535"`
		modifiedBy    string `validate:"min=1,max=100"`
		state         int    `validate:"oneof=0 1"`
		coverImageUrl string `validate:"min=1,max=255"`
	}

	var need needValid

	need.id = com.StrTo(c.Param("id")).MustInt()
	need.tagID = com.StrTo(c.Query("tag_id")).MustInt()
	need.title = c.Query("title")
	need.desc = c.Query("content")
	need.modifiedBy = c.Query("modified_by")
	need.state = com.StrTo(c.Query("state")).MustInt()
	need.coverImageUrl = c.Query("cover_image_url")

	appG := app.Gin{C: c}

	if err := validate.Struct(need); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := articleSrv.Article{
		ID:            need.id,
		TagID:         need.tagID,
		Title:         need.title,
		Desc:          need.desc,
		Content:       need.content,
		CoverImageUrl: need.coverImageUrl,
		ModifiedBy:    need.modifiedBy,
		State:         need.state,
	}

	if err := articleService.ExistByID(); err != nil {
		log.Logger.Error("article not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistArticle, nil)
		return
	}

	tagService := tagSrv.Tag{ID: need.tagID}
	if err := tagService.ExistByID(); err != nil {
		log.Logger.Error("tag not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistTag, nil)
		return
	}

	if err := articleService.Edit(); err != nil {
		log.Logger.Error("edit article fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.EditArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// DeleteArticle
// @Summary Delete a blogArticle
// @Description 删除文章
// @Tags 文章
// @Accept json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{C: c}

	if err := validate.Var(id, "min=1"); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := articleSrv.Article{ID: id}
	if err := articleService.ExistByID(); err != nil {
		log.Logger.Error("article not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistArticle, nil)
		return
	}

	if err := articleService.Delete(); err != nil {
		log.Logger.Error("delete article fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.DeleteArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
