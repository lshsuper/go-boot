package web

import (
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller interface {
	//Before Action前置操作（可以做相关逻辑：如跨域等处理）
	Before(ctx Ctx) error
	//After Action后置操作
	After(ctx Ctx) error
	//Prefix Action接口前缀（分组名）
	Prefix() string
	//Ignore 忽略的Action
	Ignore() []string
	//Recover 异常抓取操作
	Recover(ctx Ctx, err interface{})

	//===========如果需要自定义基类控制器及上下文,需要重写以下两个方法==========

	//InitCtx 初始化构建上下文策略
	InitCtx(ctx *gin.Context) Ctx
	//CallAct 执行action
	CallAct(fn interface{}, ctx Ctx)
}

//BaseController 控制器基础实现
type BaseController struct {
}

func (b *BaseController) InitCtx(ctx *gin.Context) Ctx {
	return NewBaseCtx(ctx)
}

func (b *BaseController) Ignore() (lst []string) {
	return []string{"Before", "After", "Prefix", "Ignore", "Recover", "InitCtx", "CallAct"}
}

func (b *BaseController) CallAct(fn interface{}, ctx Ctx) {
	act := fn.(func(baseCtx *BaseCtx))
	act(ctx.(*BaseCtx))
}
func (b *BaseController) Prefix() string {
	return ""
}
func (b *BaseController) Before(ctx Ctx) error {

	return nil
}
func (b *BaseController) After(ctx Ctx) error {
	return nil
}

//Recover 错误拦截
func (b *BaseController) Recover(ctx Ctx, err interface{}) {

	baseCtx := ctx.(*BaseCtx)
	if baseCtx.IsAjax() {

		//ajax请求
		baseCtx.GinCtx().JSON(http.StatusInternalServerError, gin.H{
			"message": utils.ErrorToString(err),
		})
		return

	}

	baseCtx.GinCtx().String(http.StatusInternalServerError, "", utils.ErrorToString(err))
	return

}
