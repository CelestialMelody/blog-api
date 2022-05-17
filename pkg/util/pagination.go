package util

import (
	"gin-gorm-practice/conf"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage get page number. 分页
func GetPage(c *gin.Context) int {
	result := 0
	// com.StrTo Convert string to specify type.
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * conf.AppSetting.PageSize
	}
	return result
}
