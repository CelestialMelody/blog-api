package v1

import (
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/models/blogTag"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

// GetTags 获取所有标签
// @Summary 获取多个文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = blogTag.GetTags(util.GetPage(c), setting.PageSize, maps) // maps: name state
	data["total"] = blogTag.GetTagTotal(maps)
	logging.Debug("GetTags: ", data)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// AddTags - 添加多个文章标签 POST
// @Summary 添加多个文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query string true "CreatedBy"
// @Success 200 {object} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		zap.L().Debug("AddTags: ", zap.Any("valid", valid.HasErrors()), zap.Any("valid.Errors", valid.Errors))
		if !blogTag.ExistTagByName(name) { // 判断是否存在
			if blogTag.AddTag(name, state, createdBy) { // 添加
				code = e.SUCCESS
			} else {
				logging.Info("添加多个文章标签失败")
				code = e.ERROR_ADD_TAG // 添加失败
			}
		} else {
			code = e.ERROR_EXIST_TAG           // 已存在
			for _, err := range valid.Errors { // demo 测试自己的日志
				logging.Info(err.Key, err.Message)
			}
		}
	} else {
		for _, err := range valid.Errors { // demo 测试自己的日志
			logging.Info(err.Key, err.Message)
			logging.LoggoZap.Error(
				"AddTags: ",
				zap.Any("err", err.Key),
				zap.Any("err", err.Message),
			)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

// EditTags - 编辑多个文章标签 PUT update
// @Summary 编辑多个文章标签
// @Produce  json
// @Param id query int true "ID"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {object} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTags(c *gin.Context) {
	// 获取参数
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	// 验证参数
	valid := validation.Validation{}
	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	// 验证参数
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if blogTag.ExistTagByID(id) {
			code = e.SUCCESS
			// 与AddTag 写法不同 也可以类似于AddTag那样写
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			blogTag.EditTag(id, data) //
		} else {
			code = e.ERROR_NOT_EXIST_TAG

		}
	} else {
		for _, err := range valid.Errors { // demo 测试自己的日志
			logging.Info(err.Key, err.Message)
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// DeleteTags - 删除多个文章标签
// @Summary 删除多个文章标签
// @Produce  json
// @Param id query int true "ID"
// @Success 200 {object} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTags(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")

	//binding.Validator = new(blogTag.Tag)
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if blogTag.ExistTagByID(id) {
			code = e.SUCCESS
			blogTag.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			for _, err := range valid.Errors { // demo 测试自己的日志
				logging.Info(err.Key, err.Message)
			}
		}
	} else {
		for _, err := range valid.Errors { // demo 测试自己的日志
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 打印不出来 这种方式没法直接打印
//zap.L().Debug("参数错误", zap.Any("err", valid.Errors))

// 打印不出来 测试
//logger := zap.NewExample()
//logger.Info("AddTags: ", zap.Any("err", err.Key), zap.Any("err", err.Message))
//zap.S().Error("AddTags: ", zap.Any("err", err.Key), zap.Any("err", err.Message))
