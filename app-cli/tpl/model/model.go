package model

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
	"strings"
)

//go:embed model.tpl
var tplConst string

type ModelTpl struct {
	base.BTpl
}

func NewModelTpl() *ModelTpl {
	t := &ModelTpl{}
	return t
}

type ModelArg struct {
	//MName  string          //模型名称
	Fields []ModelArgField //字段列表
	Descr  string          //描述
	TName  string          //表名称
	PName  string
}

func (ModelArg) Tag() string {
	return "modal-arg"
}

type ModelArgField struct {
	//Name   string
	//JName  string
	DType  string
	Column string
	Descr  string
}

func (tpl *ModelTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(ModelArg)

	argMap := make(map[string]interface{}, 0)
	argFiles := make([]map[string]interface{}, 0)

	if len(arg.PName) <= 0 {
		arg.PName = "model"
	}

	//表基本信息
	argMap["MName"] = strings.ReplaceAll(strings.Title(arg.TName), "_", "")
	argMap["Descr"] = arg.Descr
	argMap["TName"] = arg.TName
	argMap["PName"] = arg.PName

	//格式化属性名(执行格式化策略)
	for _, v := range arg.Fields {
		argFiles = append(argFiles, map[string]interface{}{
			"Name":   strings.ReplaceAll(strings.Title(v.Column), "_", ""),
			"JName":  v.Column,
			"Column": v.Column,
			"DType":  v.DType,
			"Descr":  v.Descr,
		})
	}
	argMap["Fields"] = argFiles

	return argMap

}

//Print 打印模板
func (tpl *ModelTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl
}

//Input 输出模板
func (tpl *ModelTpl) Input(inputPath string) base.Tpl {
	arg := tpl.GArg().(ModelArg)
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/%s_model.go", inputPath, strings.ToLower(arg.TName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl

}
