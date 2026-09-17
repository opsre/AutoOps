package model

// ExtraData 额外数据
type ExtraData map[string]interface{}

// Get 获取额外数据
func (ed ExtraData) Get(key string) interface{} {
	if ed == nil {
		return nil
	}
	return ed[key]
}

// Set 设置额外数据
func (ed ExtraData) Set(key string, value interface{}) {
	if ed == nil {
		ed = make(map[string]interface{})
	}
	ed[key] = value
}

// GetExtraString 获取额外字符串数据
func (ed ExtraData) GetExtraString(key string) string {
	if ed == nil {
		return ""
	}
	if v, ok := ed[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
