// 时间工具类
// @author xiaoRui

package util

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type HTime struct {
	time.Time
}

var (
	formatTime = "2006/01/02 15:04"
)

func (t HTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(formatTime))
	return []byte(formatted), nil
}

func (t *HTime) UnmarshalJSON(data []byte) (err error) {
	// 移除引号
	str := string(data)
	if len(str) > 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// 尝试多种时间格式解析
	formats := []string{
		"2006/01/02 15:04",              // 新格式
		"2006-01-02 15:04:05",           // 原始格式
		time.RFC3339,                     // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05",           // ISO8601 without timezone
		"2006-01-02",                     // 日期格式
	}

	var now time.Time
	for _, format := range formats {
		now, err = time.ParseInLocation(format, str, time.Local)
		if err == nil {
			*t = HTime{Time: now}
			return nil
		}
	}

	// 如果所有格式都失败，返回最后一个错误
	return err
}

func (t HTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *HTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = HTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
