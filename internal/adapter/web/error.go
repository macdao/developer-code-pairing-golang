package web

import "fmt"

// UnauthorizedError 未授权错误（Web 层使用）
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: message,
	}
}
