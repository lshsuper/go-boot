package controller

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
	"strings"
)

//go:embed controller.tpl
var tplConst string

type ControllerTpl struct {
	base.BTpl
}

func NewControllerTpl() *ControllerTpl {
	return &ControllerTpl{}
}

type ControllerArg struct {
	PkgName  string
	CtrlName string
}

func (ControllerArg) Tag() string {
	return "ctrl-arg"
}
func (tpl *ControllerTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(ControllerArg)

	argMap := make(map[string]interface{}, 0)
	if len(arg.PkgName) <= 0 {
		arg.PkgName = "controllers"
	}
	//表基本信息
	argMap["PkgName"] = arg.PkgName
	argMap["CtrlName"] = arg.CtrlName
	return argMap

}

// Print 打印模板
func (tpl *ControllerTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl
}

// Input 输出模板
func (tpl *ControllerTpl) Input(inputPath string) base.Tpl {
	arg := tpl.GArg().(ControllerArg)
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/%s_controller.go", inputPath, strings.ToLower(arg.CtrlName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl

}
