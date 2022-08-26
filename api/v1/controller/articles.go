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
	"github.com/gin-gonic/gin/binding"
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
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticleLists(c *gin.Context) {
	// 验证参数
	appG := app.Gin{C: c}
	req := com.StrTo(c.Query("tag_id")).MustInt()

	if err := validate.Var(req, "min=1"); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := articleSrv.Article{
		TagID:    req,
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	// 获取参数
	var req articleSrv.AddReq
	appG := app.Gin{C: c}

	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagService := tagSrv.Tag{ID: req.TagID}
	if err := tagService.ExistByID(); err != nil {
		log.Logger.Error("tag not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistTag, nil)
		return
	}

	articleService := articleSrv.Article{
		TagID:         req.TagID,
		Title:         req.Title,
		Desc:          req.Desc,
		Content:       req.Content,
		CoverImageUrl: req.CoverImageUrl,
		CreatedBy:     req.CreatedBy,
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
	var req articleSrv.EditReq
	appG := app.Gin{C: c}

	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := articleSrv.Article{
		ID:            req.ID,
		TagID:         req.TagID,
		Title:         req.Title,
		Desc:          req.Desc,
		Content:       req.Content,
		CoverImageUrl: req.CoverImageUrl,
		ModifiedBy:    req.ModifiedBy,
	}

	if err := articleService.ExistByID(); err != nil {
		log.Logger.Error("article not exist", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistArticle, nil)
		return
	}

	tagService := tagSrv.Tag{ID: req.TagID}
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
