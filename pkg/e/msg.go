package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	ExistTag      = 1001
	AddTag        = 1002
	NotExistTag   = 1003
	GetTagsFail   = 1004
	CountTagFail  = 1005
	AddTagFail    = 1006
	EditTagFail   = 1007
	DeleteTagFail = 1008

	NotExistArticle     = 2001
	GetArticleFail      = 2002
	GetArticleListFail  = 2003
	GetArticleCountFail = 2004
	AddArticleFail      = 2005
	EditArticleFail     = 2006
	DeleteArticleFail   = 2007

	CheckTokenFail           = 3001
	CheckAccessTokenTimeout  = 3002
	RegisterFail             = 3003
	UsernameExist            = 3004
	GenerateTokenFail        = 3005
	GenerateRefreshTokenFail = 3006
	UserNotExist             = 3007
	PasswordError            = 3008
	LoginFail                = 3009
	TokenEmpty               = 3010
	CheckRefreshTokenTimeout = 3011

	ReGenerateTokenSuccess = 3012

	UploadImageFail        = 4001
	UploadCheckImageFail   = 4003
	UploadCheckImageFormat = 4004
	UploadCheckImageSize   = 4005

	ResetRequestFail    = 5001
	BackendLoginFail    = 5002
	BackendLoginSuccess = 5003
)

const (
	CacheArticle = "ARTICLE"
	CacheTag     = "TAG"
)

var MsgFlags = map[int]string{
	Success:       "OK",
	Error:         "FAIL",
	InvalidParams: "请求参数错误",

	ExistTag:      "已存在该标签名称",
	NotExistTag:   "该标签不存在",
	AddTag:        "新增标签失败",
	GetTagsFail:   "获取多个标签失败",
	CountTagFail:  "统计标签失败",
	AddTagFail:    "新增标签失败",
	EditTagFail:   "编辑标签失败",
	DeleteTagFail: "删除标签失败",

	NotExistArticle:     "该文章不存在",
	GetArticleFail:      "获取文章失败",
	GetArticleListFail:  "获取文章列表失败",
	GetArticleCountFail: "获取文章总数失败",
	AddArticleFail:      "新增文章失败",
	EditArticleFail:     "编辑文章失败",
	DeleteArticleFail:   "删除文章失败",

	CheckTokenFail:           "Token鉴权失败",
	CheckAccessTokenTimeout:  "普通Token已超时",
	RegisterFail:             "注册失败",
	UsernameExist:            "用户名已存在",
	GenerateTokenFail:        "生成Token失败",
	GenerateRefreshTokenFail: "生成RefreshToken失败",
	UserNotExist:             "用户不存在",
	PasswordError:            "密码错误",
	LoginFail:                "登录失败",
	TokenEmpty:               "Token为空",
	CheckRefreshTokenTimeout: "RefreshToken已超时",

	ReGenerateTokenSuccess: "重新生成token成功",

	UploadImageFail:        "上传图片失败",
	UploadCheckImageFail:   "检查图片失败",
	UploadCheckImageFormat: "校验图片格式失败",
	UploadCheckImageSize:   "校验图片大小失败",

	ResetRequestFail:    "重置请求失败",
	BackendLoginFail:    "重启请求失败",
	BackendLoginSuccess: "后台登录成功",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}
