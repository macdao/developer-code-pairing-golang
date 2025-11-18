package application

import "fmt"

// ValidationError 验证错误（应用层使用）
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// InternalError 内部错误（应用层使用）
type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *InternalError {
	return &InternalError{
		Message: message,
		Err:     err,
	}
}

// NotFoundError 未找到错误（应用层使用）
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}
