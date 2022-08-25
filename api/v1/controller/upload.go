package controller

import (
	"blog-api/pkg/app"
	"blog-api/pkg/e"
	"blog-api/pkg/log"
	"blog-api/pkg/upload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
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
	appG := app.Gin{C: c}

	file, fileHeader, err := c.Request.FormFile("image") // 获取上传文件
	if err != nil {
		log.Logger.Error("get image file error", zap.Error(err))
		appG.Response(http.StatusBadRequest, e.UploadImageFail, nil)
		return
	}

	if fileHeader == nil {
		log.Logger.Error("image is nil")
		appG.Response(http.StatusBadRequest, e.UploadImageFail, nil)
		return
	}

	savePath := upload.GetImageSavePath()
	imageName := upload.GetImageName(fileHeader.Filename)

	fileUpload(imageName, file, fileHeader, c)

	appG.Response(http.StatusOK, e.Success, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
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
	code := e.Success
	data := make(map[string]interface{})
	appG := app.Gin{C: c}

	//forms := c.Request.MultipartForm // wrong; nil pointer panic
	//files, ok := forms.File["images"]
	forms, _ := c.MultipartForm()
	files, ok := forms.File["images"]
	if !ok {
		code = e.InvalidParams
		log.Logger.Error("image is nil")
		appG.Response(http.StatusBadRequest, code, nil)
	}

	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			log.Logger.Error("open image file error", zap.Error(err))
			appG.Response(http.StatusBadRequest, e.UploadImageFail, nil)
		}
		imageName := upload.GetImageName(files[i].Filename)
		savePath := upload.GetImageSavePath()

		fileUpload(imageName, file, files[i], c)

		data["image_url_"+strconv.Itoa(i)] = upload.GetImageFullUrl(imageName)
		data["image_save_url_"+strconv.Itoa(i)] = savePath + imageName
	}

	appG.Response(http.StatusOK, code, data)
}

func fileUpload(imageName string, file multipart.File, fileHeader *multipart.FileHeader, c *gin.Context) {
	fullPath := upload.GetImageFullPath()
	src := fullPath + imageName
	appG := app.Gin{C: c}

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		if !upload.CheckImageExt(imageName) {
			log.Logger.Error("image format is invalid")
			appG.Response(http.StatusBadRequest, e.UploadCheckImageFormat, nil)
		}
		if !upload.CheckImageSize(file) {
			log.Logger.Error("image size is invalid")
			appG.Response(http.StatusBadRequest, e.UploadCheckImageSize, nil)
		}
		return
	}

	if err := upload.CheckImage(fullPath); err != nil {
		log.Logger.Error("check image error", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.UploadCheckImageFail, nil)
		return
	}
	if err := c.SaveUploadedFile(fileHeader, src); err != nil { // 保存文件
		log.Logger.Error("save image error", zap.Error(err))
		appG.Response(http.StatusInternalServerError, e.UploadImageFail, nil)
	}
}
