package app

import (
	"gin-gorm-practice/pkg/log"
	"go.uber.org/zap"
)

func MarkError(err error) {
	log.Logger.Error("message:", zap.Error(err))
}
