package consts

// 通用名词
const (
	Admin            = "管理员"
	User             = "用户"
	Username         = "用户名"
	Password         = "密码"
	Phone            = "手机号"
	Email            = "邮箱"
	VerificationCode = "验证码"
	Token            = "令牌"
)

// 操作动词
const (
	Register  = "注册"
	Login     = "登录"
	Logout    = "登出"
	Operation = "操作"
	Add       = "添加"
	Update    = "更新"
	Delete    = "删除"
	Reset     = "重置"
	Invalid   = "无效"
)

// 业务实体
const (
	Artist   = "歌手"
	Song     = "歌曲"
	Playlist = "歌单"
)

// 结果状态
const (
	Success        = "成功"
	Failed         = "失败"
	Error          = "错误"
	NotFound       = "NOT_FOUND"
	DataNotFound   = "未找到相关数据"
	NotExist       = "不存在"
	AlreadyExists  = "已存在"
	NotNull        = "不能为空"
	FormatError    = "格式不正确"
	WordLimitError = "字数超出限制"
	UnknownError   = "未知错误"
	InsertFailed   = "添加失败"
	UpdateFailed   = "更新失败"
	DeleteFailed   = "删除失败"
)

// 权限/会话
const (
	AccountLocked  = "账号被锁定"
	NoPermission   = "您没有权限访问此资源"
	NotLogin       = "未登录，请先登录"
	SessionExpired = "会话过期，请重新登录"
)

// 密码相关
const (
	OldPasswordError = "原密码填写不正确"
	NewPasswordError = "新密码不能与原密码相同"
	PasswordNotMatch = "两次填写的新密码不一样"
)

// 邮件
const (
	EmailSendSuccess = "邮件发送成功"
	EmailSendFailed  = "邮件发送失败"
)

// 状态
const (
	UserStatusInvalid   = "用户状态无效"
	BannerStatusInvalid = "轮播图状态无效"
)

// 文件
const (
	FileUpload = "文件上传"
)

// 其他
const (
	InternalError = "系统内部错误"
	InvalidParams = "参数无效"
)

const (
	AdminRole = "ROLE_ADMIN"
	UserRole  = "ROLE_USER"
)
