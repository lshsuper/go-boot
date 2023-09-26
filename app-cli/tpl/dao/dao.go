package dao

import (
	_ "embed"
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"os"
	"strings"
)

//const (
//	tplConst = `package {{.PkgName}}
//import(
//     model"{{.ModelPkg}}"
//)
//
////{{.DaoName}}Dao 数据库操作
//type  {{.DaoName}}Dao struct{
//
//}
//
////New{{.DaoName}}Dao 构造器
//func New{{.DaoName}}Dao()*{{.DaoName}}Dao{
//
//     return &{{.DaoName}}Dao{}
//
//}
//
////Find 根据主键id获取数据元素
//func (dao*{{.DaoName}}Dao)Find(db *gorm.DB,id int)(m *model.{{.DaoName}}Model){
//
//     m=&model.{{.DaoName}}Model{}
//     db.Find(&m,"id=?",id)
//     if m.ID<=0{
//        m=nil
//        return
//     }
//     return
//
//}
//
//
//`
//)

//go:embed dao.tpl
var tplConst string

type DaoTpl struct {
	base.BTpl
}

func NewDaoTpl() *DaoTpl {
	return &DaoTpl{}
}

type DaoArg struct {
	PkgName  string
	DaoName  string
	ModelPkg string
}

func (arg DaoArg) Tag() string {
	return "DaoArg"
}

func (tpl *DaoTpl) FArgs() map[string]interface{} {

	arg := tpl.GArg().(DaoArg)

	argMap := make(map[string]interface{}, 0)
	if len(arg.PkgName) <= 0 {
		arg.PkgName = "dao"
	}
	//表基本信息
	argMap["PkgName"] = arg.PkgName
	argMap["DaoName"] = arg.DaoName
	argMap["ModelPkg"] = arg.ModelPkg

	return argMap

}

//Print 打印模板
func (tpl *DaoTpl) Print() base.Tpl {

	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	_ = t.Execute(os.Stdout, tpl.FArgs())
	return tpl

}

//Input 输出模板
func (tpl *DaoTpl) Input(inputPath string) base.Tpl {
	arg := tpl.GArg().(DaoArg)
	t, err := tpl.GTplInstance(tplConst)
	if err != nil {
		panic(err.Error())
	}
	fs, err := os.OpenFile(fmt.Sprintf("%s/%s_dao.go", inputPath, strings.ToLower(arg.DaoName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	_ = t.Execute(fs, tpl.FArgs())
	return tpl
}
