package util

import (
	"gin-gorm-practice/conf/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage get page number. 分页; golang写的一个分页控件
func GetPage(c *gin.Context) int {
	result := 0
	// com.StrTo Convert string to specify type.
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result
}

func GetPageSize(c *gin.Context) int {
	result := 0
	pageSize, _ := com.StrTo(c.Query("pageSize")).Int()
	if pageSize > 0 {
		result = pageSize
	}
	return result
}
