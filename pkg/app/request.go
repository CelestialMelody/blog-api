package app

import (
	"gin-gorm-practice/pkg/logging"
	"github.com/beego/beego/v2/core/validation"
	"go.uber.org/zap"
)

func MarKError(err interface{}) {
	switch err.(type) {
	case error:
		e := err.(error)
		logging.LoggoZap.Error("validate error", zap.Error(e))
	case *validation.Error:
		e := err.(*validation.Error)
		logging.LoggoZap.Error(e.Key, zap.Any("massage:", e.Value))
	}
}
