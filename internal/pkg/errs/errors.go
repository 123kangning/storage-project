package errs

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewErrorCode 创建新的错误码
func NewErrorCode(code int, message string) *ErrorCode {
	return &ErrorCode{
		Code:    code,
		Message: message,
	}
}

// Error 实现 error 接口
func (e *ErrorCode) Error() string {
	return e.Message
}
