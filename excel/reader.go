package excel

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
	"strings"
)

type readerProvider struct {
	fs *excelize.File
}

func NewReaderProvider(fileName string) *readerProvider {
	fs, _ := excelize.OpenFile(fileName)
	return &readerProvider{fs: fs}
}

//ReadMap 读成map结构
func (e *readerProvider) ReadMap() (res []map[string]string, err error) {

	sheetName := e.fs.GetSheetName(0)
	res, err = e.ReadMapBySheet(sheetName)
	return

}

//ReadMapBySheet 读成map结构(by sheet)
func (e *readerProvider) ReadMapBySheet(sheetName string) (res []map[string]string, err error) {

	rows, err := e.fs.GetRows(sheetName)
	if err != nil {
		return
	}
	res = make([]map[string]string, 0)
	h := rows[0]
	for index, row := range rows {
		if index == 0 || row == nil {
			continue
		}
		curRow := map[string]string{}
		for k, v := range h {
			if (k + 1) > len(row) {
				curRow[v] = ""
			} else {
				curRow[v] = row[k]
			}

		}
		res = append(res, curRow)
	}

	return

}

//Close 释放资源
func (e *readerProvider) Close() {

	_ = e.fs.Close()
	e.fs = nil
	return

}

//ReadModel 映射模型的读
func (e *readerProvider) ReadModel(res interface{}) error {

	vt := reflect.ValueOf(res)
	vtt := vt.Elem()
	if vtt.Type().Kind() != reflect.Slice {
		return errors.New("res must slice")
	}
	mt := vtt.Type().Elem()

	//读取excel
	rows, err := e.SimpleReader()
	if err != nil {
		return err
	}
	if len(rows) <= 0 {
		return nil
	}
	hMap := make(map[string]int)
	for hk, h := range rows[0] {
		hMap[strings.ReplaceAll(h, " ", "")] = hk
	}

	//解析除第一行外的所有行数据
	for rk, r := range rows {
		if rk == 0 {
			continue
		}

		//组装结构体
		val := reflect.New(mt).Elem()
		//拼装行元素
		for i := 0; i < val.NumField(); i++ {

			var (
				curField   = val.Field(i)
				colName    = mt.Field(i).Tag.Get("x-col")
				cIndex, ok = hMap[colName]
			)

			if !ok {
				continue
			}

			switch curField.Type().Kind() {

			case reflect.String:
				curField.SetString(r[cIndex])
			case reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int, reflect.Int16:
				fval, _ := strconv.Atoi(r[cIndex])
				curField.SetInt(int64(fval))
			case reflect.Float64, reflect.Float32:
				fval, _ := strconv.ParseFloat(r[cIndex], 64)
				curField.SetFloat(fval)
			case reflect.Bool:
				fval, _ := strconv.ParseBool(r[cIndex])
				curField.SetBool(fval)
			}
		}
		//补充返回切片列表
		vtt = reflect.Append(vtt, val)

	}
	//构造指针切片
	vt.Elem().Set(vtt)
	return nil

}

//SimpleReader 简单读
func (e *readerProvider) SimpleReader() ([][]string, error) {
	sheetName := e.fs.GetSheetName(0)
	return e.Reader(sheetName)
}

//Reader 读
func (e *readerProvider) Reader(sheetName string) ([][]string, error) {
	return e.fs.GetRows(sheetName)
}

//ForEach 遍历行
func (e *readerProvider) ForEach(fn func(k int, r []string)) error {

	rows, err := e.SimpleReader()
	if err != nil {
		return err
	}

	for k, v := range rows {
		fn(k, v)
	}

	return nil

}
