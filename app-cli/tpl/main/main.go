package _main

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
)

//go:embed main.tpl
var tplConst string

type MainTpl struct {
	base.BTpl
}

func NewMainTpl() *MainTpl {
	return &MainTpl{}
}

type MainArg struct {
	PName string
}

func (MainArg) Tag() string {
	return "main-arg"
}
func (tpl *MainTpl) FArgs() map[string]interface{} {

	//arg := tpl.GArg().(MainArg)

	argMap := make(map[string]interface{}, 0)

	return argMap

}

//Print 打印模板
func (tpl *MainTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl
}

//Input 输出模板
func (tpl *MainTpl) Input(inputPath string) base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/main.go", inputPath), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl

}
