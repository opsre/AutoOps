package collx

// M 是map[string]any的别名
type M map[string]any

// Kvs 键值对切片
type Kvs []Kv

type Kv struct {
	Key   string
	Value any
}

// NewKvs 创建键值对切片
func NewKvs(args ...any) Kvs {
	if len(args)%2 != 0 {
		return Kvs{}
	}
	result := make(Kvs, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		result = append(result, Kv{
			Key:   args[i].(string),
			Value: args[i+1],
		})
	}
	return result
}

// ArrayMap 数组映射
func ArrayMap[T any, R any](arr []T, fn func(T) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

// ArrayFilter 数组过滤
func ArrayFilter[T any](arr []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range arr {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// ArrayToMap 数组转Map
func ArrayToMap[T any, K comparable](arr []T, keyFn func(T) K) map[K]T {
	result := make(map[K]T)
	for _, v := range arr {
		result[keyFn(v)] = v
	}
	return result
}

// AsArray 转换为数组
func AsArray[T any](v any) []T {
	if v == nil {
		return []T{}
	}
	if arr, ok := v.([]T); ok {
		return arr
	}
	return []T{}
}

// ArrayContains 数组是否包含
func ArrayContains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// ArrayAnyMatches 数组任意匹配
func ArrayAnyMatches[T any](arr []T, fn func(T) bool) bool {
	for _, v := range arr {
		if fn(v) {
			return true
		}
	}
	return false
}

// AnyMatch 任意匹配(别名)
func AnyMatch[T any](arr []T, fn func(T) bool) bool {
	return ArrayAnyMatches(arr, fn)
}

// ArrayRemoveFunc 删除满足条件的元素
func ArrayRemoveFunc[T any](arr []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range arr {
		if !fn(v) {
			result = append(result, v)
		}
	}
	return result
}
