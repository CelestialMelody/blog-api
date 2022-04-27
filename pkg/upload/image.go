package upload

import (
	"fmt"
	"gin-gorm-practice/conf/setting"
	"gin-gorm-practice/pkg/file"
	"gin-gorm-practice/pkg/logging"
	"gin-gorm-practice/pkg/util"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImageSavePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImagePrefixUrl() string {
	return setting.AppSetting.ImagePrefixUrl
}

func GetRuntimeRootPath() string {
	return setting.AppSetting.RuntimeRootPath
}

func GetImageFullUrl(name string) string {
	return GetImagePrefixUrl() + "/" + GetImageSavePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)                     // 获取文件后缀
	fileName := strings.TrimSuffix(name, ext) // 去除文件名后缀
	//fileName = strings.Replace(fileName, "-", "_", -1) // 替换文件名中的 - 为 _
	fileName = util.EncodeMD5(fileName) + ext // 加密文件名
	return fileName
}

func GetImageFullPath() string {
	return GetRuntimeRootPath() + GetImageSavePath()
}

// CheckImageExt 检查图片后缀
func CheckImageExt(fileName string) bool {
	ext := path.Ext(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExt {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool { // 检查图片大小
	size, err := file.GetSize(f)
	if err != nil {
		logging.LoggoZap.Warn("获取文件大小失败", zap.Error(err))
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd() // 获取当前目录
	if err != nil {
		logging.LoggoZap.Error("os.Getwd err", zap.Error(err))
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		logging.LoggoZap.Error("file.IsNotExistMkDir err", zap.Error(err))
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}
	perm := file.CheckPermission(src)
	if perm == true {
		logging.LoggoZap.Error("file.CheckPermission err", zap.Error(err))
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
