package app

import (
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
	return
}

func MarkError(err error) {
	log.Logger.Error("err_msg:", zap.Error(err))
}
