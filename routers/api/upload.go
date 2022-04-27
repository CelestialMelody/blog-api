package api

import (
	"fmt"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/upload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	file, image, err := c.Request.FormFile("image") // 获取上传文件
	if err != nil {
		logging.LoggoZap.Error("Failed to get image", zap.Any("err", err))
		code = e.ERROR
		c.JSON(200, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
		err = fmt.Errorf("image is nil")
		logging.LoggoZap.Error("Invalid params", zap.Any("err", err))
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImageSavePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
			if !upload.CheckImageExt(imageName) {
				err = fmt.Errorf("image format is invalid")
				logging.LoggoZap.Error("Image format is invalid", zap.Any("err", err))
			}
			if !upload.CheckImageSize(file) {
				err = fmt.Errorf("image size is invalid")
				logging.LoggoZap.Error("Image size is invalid", zap.Any("err", err))
			}
			//logging.LoggoZap.Error("Check image format failed", zap.Any("err", err))
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				logging.LoggoZap.Error("Check image failed", zap.Any("err", err))
			} else if err := c.SaveUploadedFile(image, src); err != nil { // 保存文件
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
				logging.LoggoZap.Error("Save image failed", zap.Any("err", err))
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
