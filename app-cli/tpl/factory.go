package tpl

import (
	"gitee.com/lshsuper/go-boot/app-cli/tpl/base"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/controller"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/dao"
	_main "gitee.com/lshsuper/go-boot/app-cli/tpl/main"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/mod"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/model"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/router"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/service"
)

type TType string

const (
	ModelTpl  TType = "Model"
	SrvTpl    TType = "Srv"
	ModTpl    TType = "Mod"
	CtrlTpl   TType = "Ctrl"
	RouterTpl TType = "Router"
	MainTpl   TType = "Main"
	DaoTpl    TType = "Dao"
)

func (e TType) Tpl() base.Tpl {

	switch e {
	case ModelTpl:
		return model.NewModelTpl()
	case ModTpl:
		return mod.NewModTpl()
	case CtrlTpl:
		return controller.NewControllerTpl()
	case SrvTpl:
		return service.NewServiceTpl()
	case RouterTpl:
		return router.NewRouterTpl()
	case MainTpl:
		return _main.NewMainTpl()
	case DaoTpl:
		return dao.NewDaoTpl()
	default:
		panic("不存在的模板类型")
	}

}

//func ModelTpl() *model.ModelTpl {
//	return model.NewModelTpl()
//}
