package api

import (
	"fmt"
	"gin-gorm-practice/pkg/e"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/upload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// UploadImage upload image
// @Summary Upload image
// @Tags 图片上传
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "{"code":0,"message":"ok","data":{}}"
// @Router /upload [post]
func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	file, fileHeader, err := c.Request.FormFile("image") // 获取上传文件
	logging.LoggoZap.Info("file fileHeader", zap.Any("file", file), zap.Any("fileHeader", fileHeader))
	if err != nil {
		logging.LoggoZap.Error("Failed to get image", zap.Any("err", err))
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if fileHeader == nil {
		code = e.INVALID_PARAMS
		err = fmt.Errorf("image is nil")
		logging.LoggoZap.Error("Invalid params", zap.Any("err", err))
	} else {
		imageName := upload.GetImageName(fileHeader.Filename)
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
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				logging.LoggoZap.Error("Check image failed", zap.Any("err", err))
			} else if err := c.SaveUploadedFile(fileHeader, src); err != nil { // 保存文件
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

// UploadImages UploadFile upload file
// @Summary Upload file
// @Tags 文件上传
// @Produce json
// @Param file formData file true "File"
// @Success 200 {string} string "{"code":0,"message":"ok","data":{}}"
// @Router /upload/file [post]
func UploadImages(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	//forms := c.Request.MultipartForm
	//files, ok := forms.File["images"]
	forms, _ := c.MultipartForm()
	files, ok := forms.File["images"]
	if !ok {
		code = e.INVALID_PARAMS
		err := fmt.Errorf("images is nil")
		logging.LoggoZap.Error("Invalid params", zap.Any("err", err))
	}
	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			code = e.ERROR
			err = fmt.Errorf("failed to get image")
			logging.LoggoZap.Error("Failed to get image", zap.Any("err", err))

			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
		}
		imageName := upload.GetImageName(files[i].Filename)
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
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				logging.LoggoZap.Error("Check image failed", zap.Any("err", err))
			} else if err := c.SaveUploadedFile(files[i], src); err != nil { // 保存文件
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
				logging.LoggoZap.Error("Save image failed", zap.Any("err", err))
			} else {
				data["image_url_"+strconv.Itoa(i)] = upload.GetImageFullUrl(imageName)
				data["image_save_url_"+strconv.Itoa(i)] = savePath + imageName

				c.JSON(http.StatusOK, gin.H{
					"code": code,
					"msg":  e.GetMsg(code),
					"data": data,
				})
			}
		}
	}
}
