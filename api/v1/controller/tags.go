package controller

import (
	"blog-api/conf"
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

// GetTagLists - 获取多个文章标签 GET
// @Summary GetTagLists
// @Produce json
// @Tags 标签
// @Description Get multiple blogArticle tags
// @Param name query string false "标签名称"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTagLists(c *gin.Context) {
	appG := app.Gin{C: c}
	req := c.Query("name")

	if err := validate.Var(req, "required,max=100"); err != nil {
		log.Logger.Error("invalid Params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagServeice := tagSrv.Tag{
		Name:     req,
		PageNum:  util.GetPage(c),
		PageSize: conf.AppConfig.PageSize,
	}

	tagLists, err := tagServeice.GetAll()
	if err != nil {
		log.Logger.Error("get tag list failed", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.GetTagsFail, nil)
	}

	total, err := tagServeice.Count()
	if err != nil {
		log.Logger.Error("get tag count failed", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.CountTagFail, nil)
	}

	appG.Response(http.StatusOK, e.Success, map[string]interface{}{
		"lists": tagLists,
		"total": total,
	})
}

// AddTags - 添加多个文章标签 POST
// @Summary AddTags
// @Tags 标签
// @Description Add multiple blogArticle tags
// @Produce json
// @Param name query string true "标签名称"
// @Param created_by query string true "创建人"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	var req tagSrv.AddReq
	appG := app.Gin{C: c}

	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid Params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagService := tagSrv.Tag{
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
	}

	if err := tagService.ExistByName(); err == nil {
		log.Logger.Error("tag has existed", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.ExistTag, nil)
		return
	}

	if err := tagService.Add(); err != nil {
		log.Logger.Error("add tag failed", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.AddTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// EditTags - 编辑多个文章标签 PUT update
// @Summary EditTags
// @Tags 标签
// @Description Edit multiple blogArticle tags
// @Produce json
// @Param id path int true "ID"
// @Param name query string true "标签名称"
// @Param modified_by query string true "修改人"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTags(c *gin.Context) {
	var req tagSrv.EditReq
	appG := app.Gin{C: c}

	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagService := tagSrv.Tag{
		ID:         req.ID,
		Name:       req.Name,
		ModifiedBy: req.ModifiedBy,
	}

	if err := tagService.ExistByID(); err != nil {
		log.Logger.Error("not exist tag", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistTag, nil)
		return
	}

	if err := tagService.Edit(); err != nil {
		log.Logger.Error("edit tag fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.EditTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}

// DeleteTags - 删除多个文章标签
// @Summary DeleteTags
// @Tags 标签
// @Produce json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTags(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{C: c}

	if err := validate.Var(id, "required,min=1"); err != nil {
		log.Logger.Error("invalid params", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	tagService := tagSrv.Tag{ID: id}
	if err := tagService.ExistByID(); err == nil {
		log.Logger.Error("not exist tag", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.NotExistTag, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		log.Logger.Error("delete tag fail", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.DeleteTagFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.Success, nil)
}
