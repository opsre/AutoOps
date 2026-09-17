package errorx

import "fmt"

// BizError 业务错误
type BizError struct {
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

// NewBiz 创建业务错误
func NewBiz(msg string) error {
	return &BizError{Message: msg}
}

// NewBizf 创建格式化业务错误
func NewBizf(format string, args ...interface{}) error {
	return &BizError{Message: fmt.Sprintf(format, args...)}
}
