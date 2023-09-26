package cmd

import (
	"gitee.com/lshsuper/go-boot/app-cli/tpl"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/controller"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/dao"
	_main "gitee.com/lshsuper/go-boot/app-cli/tpl/main"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/mod"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/model"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/router"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/service"
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/urfave/cli"
	"os"
	"path"
)

//APICommand api应用构建命令
var APICommand = cli.Command{
	Name:        "api",
	Description: "构建一个api应用",
	ShortName:   "a",
	Usage:       "构建一个api应用",
	Action: func(ctx *cli.Context) {
		var (
			n      = ctx.String("n")
			curDir = utils.GetCurrentDirectory()
		)

		_ = os.Mkdir(path.Join(curDir, n), os.ModePerm)
		_ = os.Mkdir(path.Join(curDir, n, controllerRootPkgName), os.ModePerm)
		ctrlTpl := tpl.CtrlTpl.Tpl()
		ctrlTpl.SArg(controller.ControllerArg{
			CtrlName: "Default",
		})
		ctrlTpl.Input(path.Join(curDir, n, controllerRootPkgName))

		_ = os.Mkdir(path.Join(curDir, n, modelRootPkgName), os.ModePerm)
		modelTpl := tpl.ModelTpl.Tpl()
		modelTpl.SArg(model.ModelArg{
			Fields: []model.ModelArgField{{Descr: "默认", DType: "int", Column: "ID"}},
			Descr:  "默认模型",
			TName:  "default",
		})

		modelTpl.Input(path.Join(curDir, n, modelRootPkgName))

		_ = os.Mkdir(path.Join(curDir, n, routerRootPkgName), os.ModePerm)
		routerTpl := tpl.RouterTpl.Tpl()
		routerTpl.SArg(router.RouterArg{
			CtrlPkg:  "",
			CtrlName: "Test",
		})
		routerTpl.Input(path.Join(curDir, n, routerRootPkgName))

		_ = os.Mkdir(path.Join(curDir, n, serviceRootPkgName), os.ModePerm)
		srvTpl := tpl.SrvTpl.Tpl()
		srvTpl.SArg(service.ServiceArg{
			SrvName: "Default",
		})
		srvTpl.Input(path.Join(curDir, n, serviceRootPkgName))

		_ = os.Mkdir(path.Join(curDir, n, daoRootPkgName), os.ModePerm)
		daoTpl := tpl.DaoTpl.Tpl()
		daoTpl.SArg(dao.DaoArg{
			DaoName: "Default",
		})
		daoTpl.Input(path.Join(curDir, n, daoRootPkgName))

		//初始化mod
		modTpl := tpl.ModTpl.Tpl()
		modTpl.SArg(mod.ModArg{
			PName: n,
		})
		modTpl.Input(path.Join(curDir, n))

		//初始化main
		mainTpl := tpl.MainTpl.Tpl()
		mainTpl.SArg(_main.MainArg{})
		mainTpl.Input(path.Join(curDir, n))

		//工具目录
		os.Mkdir(path.Join(curDir, n, pkgRootPkgName), os.ModePerm)

	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "n",
			Value: "app",
			Usage: "应用名称",
		},
	},
}
