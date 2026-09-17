package stringx

import (
	"bytes"
	"strings"
	"text/template"
	"github.com/google/uuid"
)

// IsEmpty 判断字符串是否为空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 判断字符串是否不为空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// DefaultIfEmpty 如果为空则返回默认值
func DefaultIfEmpty(s, defaultVal string) string {
	if IsEmpty(s) {
		return defaultVal
	}
	return s
}

// TrimSpaceAndBr 去除空格和换行
func TrimSpaceAndBr(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// TemplateParse 模板解析
func TemplateParse(name, text string, data interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RandUUID 生成随机UUID
func RandUUID() string {
	return uuid.New().String()
}
