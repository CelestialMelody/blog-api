package v1

import (
	"gin-gorm-practice/conf"
	"gin-gorm-practice/models/blogArticle"
	"gin-gorm-practice/pkg/app"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/util"
	"gin-gorm-practice/pkg/validate"
	"gin-gorm-practice/service/articleS"
	"gin-gorm-practice/service/tagS"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
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
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := blogArticle.ExistArticleByID(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}

	articleService := articleS.Article{ID: id}
	article, err := articleService.Get()
	if err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
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
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := articleS.Article{
		TagID:    need.tagId,
		State:    need.state,
		PageNum:  util.GetPage(c),
		PageSize: conf.AppSetting.PageSize,
	}

	articleLists, err := articleService.GetAll()
	if err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_LIST_FAIL, nil)
		return
	}

	total, err := articleService.Count()
	if err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_COUNT_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articleLists
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
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
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tagS.Tag{ID: need.tagId}
	if err := tagService.ExistByID(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := articleS.Article{
		TagID:         need.tagId,
		Title:         need.title,
		Desc:          need.desc,
		Content:       need.content,
		CoverImageUrl: need.coverImageUrl,
		CreatedBy:     need.createdBy,
		State:         need.state,
	}

	if err := articleService.Add(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := articleS.Article{
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
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tagS.Tag{ID: need.tagID}
	if err := tagService.ExistByID(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := articleService.Edit(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
		app.MarkError(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := articleS.Article{ID: id}
	if err := articleService.ExistByID(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	if err := articleService.Delete(); err != nil {
		app.MarkError(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
