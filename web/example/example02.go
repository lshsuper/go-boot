package example

import (
	"gitee.com/lshsuper/go-boot/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

//ApiCtx 自定义请求上下文
type ApiCtx struct {
	web.BaseCtx
}

func NewApiCtx(ctx *gin.Context) *ApiCtx {
	c := &ApiCtx{}
	c.Context = ctx
	return c
}

//ApiController 自定义基类控制器
type ApiController struct {
	web.BaseController
}

func (b *ApiController) InitCtx(ctx *gin.Context) web.Ctx {
	return NewApiCtx(ctx)
}

func (b *ApiController) CallAct(fn interface{}, ctx web.Ctx) {
	act := fn.(func(baseCtx *ApiCtx))
	act(ctx.(*ApiCtx))
}

type Example02Controller struct {
	ApiController
}

func NewExample02Controller() *Example02Controller {
	return &Example02Controller{}
}

//@Router /e/test
//@Method post
func (ctrl *Example02Controller) Test(ctx *ApiCtx) {

	ctx.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
	return

}

func Example02() {

	s := web.NewServer(web.WebConf{
		Boot: "/",
	})
	//常规method模式
	s.AutoRouter(web.NormalMethod, NewExample02Controller())
	s.Run(":10086")
}
