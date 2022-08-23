package app

import (
	"blog-api/pkg/log"
	"go.uber.org/zap"
)

func MarkError(err error) {
	log.Logger.Error("message:", zap.Error(err))
}
