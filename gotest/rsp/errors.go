package rsp

// ErrorMessages 错误码和错误信息的映射表
var ErrorMessages = map[int]string{
	// 用户相关错误
	10001: "用户已存在",    // 用户已经存在
	10002: "用户不存在",    // 用户不存在
	10003: "用户名或密码错误", // 用户名或密码不匹配

	// 数据库相关错误
	20001: "数据库连接失败", // 数据库连接失败
	20002: "数据库查询失败", // 数据库查询失败
	20003: "数据库插入失败", // 数据库插入失败
	20004: "数据库迁移失败", // 数据库迁移失败

	// 请求参数错误
	30001: "缺少请求参数", // 缺少必要的请求参数
	30002: "请求参数无效", // 请求参数格式无效

	// 系统级错误
	40001: "内部服务器错误", // 服务器内部错误
	40002: "服务不可用",   // 服务不可用

	// JWT 相关错误
	50001: "JWT 生成失败", // JWT 生成失败
	50002: "密码加密失败",   // 密码加密失败

	// 系统数据库问题
	60001: "数据库问题", // 服务器数据库问题
}

type ErrorResponse struct {
	Code    int         `json:"code"`           // 错误码
	Message string      `json:"message"`        // 错误描述
	Data    interface{} `json:"data,omitempty"` // 附加数据，通常为额外的错误信息
}

// NewErrorResponse 用于创建一个新的 ErrorResponse
func NewErrorResponse(code int, data interface{}) *ErrorResponse {
	// 查找错误信息
	message, exists := ErrorMessages[code]
	if !exists {
		message = "未知错误" // 如果没有找到对应的错误信息，使用默认信息
	}

	return &ErrorResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
