package mod

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
)

//go:embed mod.tpl
var tplConst string

type ModTpl struct {
	base.BTpl
}

func NewModTpl() *ModTpl {
	return &ModTpl{}
}

type ModArg struct {
	PName string
}

func (ModArg) Tag() string {
	return "mod-arg"
}
func (tpl *ModTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(ModArg)

	argMap := make(map[string]interface{}, 0)

	//表基本信息
	argMap["PName"] = arg.PName

	return argMap

}

//Print 打印模板
func (tpl *ModTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl
}

//Input 输出模板
func (tpl *ModTpl) Input(inputPath string) base.Tpl {
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/go.mod", inputPath), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl

}
