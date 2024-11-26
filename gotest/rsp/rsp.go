package rsp

type ErrorResponse struct {
	Code     int         `json:"code"`           // 错误码
	Message  string      `json:"message"`        // 错误描述
	Data     interface{} `json:"data,omitempty"` // 附加数据，通常为额外的错误信息
	Errinput interface{} `json:"input"`          // 附加数据，通常为额外的错误信息
}

// NewErrorResponse 用于创建一个新的 ErrorResponse
func NewErrorResponse(code int, data interface{}, input interface{}) *ErrorResponse {
	// 查找错误信息
	message, exists := ErrorMessages[code]
	if !exists {
		message = "未知错误" // 如果没有找到对应的错误信息，使用默认信息
	}

	return &ErrorResponse{
		Code:     code,
		Message:  message,
		Data:     data,
		Errinput: input,
	}
}

// NewSuccessResponse 成功返回
func NewSuccessResponse(code int, data interface{}) *ErrorResponse {
	// 查找错误信息
	message, exists := SuccessMessages[code]
	if !exists {
		message = "未定义的成功" // 如果没有找到对应的错误信息，使用默认信息
	}
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
