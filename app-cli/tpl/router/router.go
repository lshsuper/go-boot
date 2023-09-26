package router

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
)

//go:embed router.tpl
var tplConst string

type RouterTpl struct {
	base.BTpl
}

func NewRouterTpl() *RouterTpl {
	return &RouterTpl{}
}

type RouterArg struct {
	PkgName  string
	CtrlPkg  string
	CtrlName string
}

func (arg RouterArg) Tag() string {
	return "RouterArg"
}

func (tpl *RouterTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(RouterArg)
	argMap := make(map[string]interface{}, 0)
	if len(arg.PkgName) <= 0 {
		arg.PkgName = "routers"
	}
	if len(arg.CtrlPkg) <= 0 {
		arg.CtrlPkg = "controller"
	}
	//表基本信息
	argMap["PkgName"] = arg.PkgName
	argMap["CtrlPkg"] = arg.CtrlPkg
	argMap["CtrlName"] = arg.CtrlName
	return argMap

}

//Print 打印模板
func (tpl *RouterTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl
}

//Input 输出模板
func (tpl *RouterTpl) Input(inputPath string) base.Tpl {

	//arg := tpl.GArg().(RouterArg)
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/router.go", inputPath), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl
}
