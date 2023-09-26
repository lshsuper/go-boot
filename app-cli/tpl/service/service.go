package service

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
	"strings"
)

//go:embed service.tpl
var tplConst string

type ServiceTpl struct {
	base.BTpl
}

func NewServiceTpl() *ServiceTpl {
	return &ServiceTpl{}
}

type ServiceArg struct {
	PkgName string
	SrvName string
}

func (arg ServiceArg) Tag() string {
	return "ServiceArg"
}

func (tpl *ServiceTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(ServiceArg)

	argMap := make(map[string]interface{}, 0)
	if len(arg.PkgName) <= 0 {
		arg.PkgName = "services"
	}
	//表基本信息
	argMap["PkgName"] = arg.PkgName
	argMap["SrvName"] = arg.SrvName

	return argMap

}

//Print 打印模板
func (tpl *ServiceTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl

}

//Input 输出模板
func (tpl *ServiceTpl) Input(inputPath string) base.Tpl {
	arg := tpl.GArg().(ServiceArg)
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("./%s_service.go", strings.ToLower(arg.SrvName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl
}
