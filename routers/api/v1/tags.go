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
	"net/http"
)

// GetTags - 获取多个文章标签 GET
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
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

	//logrus.Println(gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
}

// EditTags - 编辑多个文章标签 PUT update
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
	//valid.Required(name, "name").Message("名称不能为空")
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
			for _, err := range valid.Errors { // demo 测试自己的日志
				logging.Info(err.Key, err.Message)
			}
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
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
