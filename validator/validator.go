package validator

import (
	"reflect"
	"strings"
)

//ValidMulti 验证多条件
func ValidMulti(key string, val interface{}, allErr bool, vts ...ValidType) (es []error) {

	for _, vt := range vts {
		if err := vt.Valid(key, val); err != nil {
			es = append(es, err)
			if !allErr {
				return
			}
		}
	}

	return

}

//ValidStruct 验证结构体
func ValidStruct(m interface{}, allErr bool) (es []error) {

	var (
		vv = reflect.ValueOf(m)
		vt = vv.Type()
	)
	for i := 0; i < vt.NumField(); i++ {

		ft := vt.Field(i)
		if tv := ft.Tag.Get(validTag); len(tv) > 0 {
			fv := vv.Field(i)

			//需要验证
			fvArr := strings.Split(tv, ";")
			for _, fvRow := range fvArr {
				if err := ValidType(fvRow).Valid(ft.Name, fv.Interface()); err != nil {
					es = append(es, err)
					if !allErr {
						return
					}
				}

			}

		}
	}

	return nil

}
