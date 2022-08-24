package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	ErrorExistTag      = 1001
	ErrorAddTag        = 1002
	ErrorNotExistTag   = 1003
	ErrorGetTagsFail   = 1004
	ErrorCountTagFail  = 1005
	ErrorAddTagFail    = 1006
	ErrorEditTagFail   = 1007
	ErrorDeleteTagFail = 1008

	ErrorNotExistArticle     = 2001
	ErrorGetArticleFail      = 2002
	ErrorGetArticleListFail  = 2003
	ErrorGetArticleCountFail = 2004
	ErrorAddArticleFail      = 2005
	ErrorEditArticleFail     = 2006
	ErrorDeleteArticleFail   = 2007

	ErrorAuthCheckTokenFail    = 3001
	ErrorAuthCheckTokenTimeout = 3002
	ErrorAuthToken             = 3003
	ErrorRegisterFail          = 3004

	ErrorUploadImageFail        = 4001
	ErrorUploadCheckImageFail   = 4003
	ErrorUploadCheckImageFormat = 4004
	ErrorUploadCheckImageSize   = 4005

	ErrorNotExistUser = 5001
)

const (
	CacheArticle = "ARTICLE"
	CacheTag     = "TAG"
)

var MsgFlags = map[int]string{
	Success:       "OK",
	Error:         "FAIL",
	InvalidParams: "请求参数错误",

	ErrorExistTag:      "已存在该标签名称",
	ErrorNotExistTag:   "该标签不存在",
	ErrorAddTag:        "新增标签失败",
	ErrorGetTagsFail:   "获取多个标签失败",
	ErrorCountTagFail:  "统计标签失败",
	ErrorAddTagFail:    "新增标签失败",
	ErrorEditTagFail:   "编辑标签失败",
	ErrorDeleteTagFail: "删除标签失败",

	ErrorNotExistArticle:     "该文章不存在",
	ErrorGetArticleFail:      "获取文章失败",
	ErrorGetArticleListFail:  "获取文章列表失败",
	ErrorGetArticleCountFail: "获取文章总数失败",
	ErrorAddArticleFail:      "新增文章失败",
	ErrorEditArticleFail:     "编辑文章失败",
	ErrorDeleteArticleFail:   "删除文章失败",

	ErrorAuthCheckTokenFail:    "Token鉴权失败",
	ErrorAuthCheckTokenTimeout: "Token已超时",
	ErrorAuthToken:             "Token生成失败",
	ErrorNotExistUser:          "该用户不存在",
	ErrorRegisterFail:          "注册失败",

	ErrorUploadImageFail:        "上传图片失败",
	ErrorUploadCheckImageFail:   "检查图片失败",
	ErrorUploadCheckImageFormat: "校验图片格式失败",
	ErrorUploadCheckImageSize:   "校验图片大小失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}
