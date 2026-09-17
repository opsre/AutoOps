package anyx

import "reflect"

// ToString 转换为字符串
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// ToInt 转换为整数
func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

// IsBlank 判断值是否为空
func IsBlank(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return val.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	default:
		return false
	}
}
