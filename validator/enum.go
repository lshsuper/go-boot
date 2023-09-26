package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type ValidType string

const (
	Phone ValidType = "phone"
	Email ValidType = "email"
	Ip    ValidType = "ip"
	Empty ValidType = "empty"
)

//Valid 单条件验证
func (e ValidType) Valid(key string, val interface{}) (err error) {

	switch e {
	case Phone:
		if ok := phoneRegex.MatchString(val.(string)); !ok {
			err = errFormat(key, "不是正确的手机号")
		}
	case Email:
		if ok := emailRegex.MatchString(val.(string)); !ok {
			err = errFormat(key, "不是正确的邮箱")
		}
	case Ip:
		if ok := ipRegex.MatchString(val.(string)); !ok {
			err = errFormat(key, "不是正确的IP地址")
		}
	case Empty:
		if ok := isEmpty(val); !ok {
			err = errFormat(key, "不能为空")
		}
	}
	return
}

//errFormat 错误信息格式化
func errFormat(key, msg string) error {
	return errors.New(fmt.Sprintf("参数【%s】%s", key, msg))
}

//isEmpty 是否为空
func isEmpty(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(strings.TrimSpace(str)) > 0
	}
	if _, ok := obj.(bool); ok {
		return true
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if i, ok := obj.(uint); ok {
		return i != 0
	}
	if i, ok := obj.(int8); ok {
		return i != 0
	}
	if i, ok := obj.(uint8); ok {
		return i != 0
	}
	if i, ok := obj.(int16); ok {
		return i != 0
	}
	if i, ok := obj.(uint16); ok {
		return i != 0
	}
	if i, ok := obj.(uint32); ok {
		return i != 0
	}
	if i, ok := obj.(int32); ok {
		return i != 0
	}
	if i, ok := obj.(int64); ok {
		return i != 0
	}
	if i, ok := obj.(uint64); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}
