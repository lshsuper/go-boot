package example

import (
	"gitee.com/lshsuper/go-boot/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestController struct {
	web.BaseController
}

func NewTestController() *TestController {
	return &TestController{}
}

func (ctrl *TestController) Test(ctx *web.BaseCtx) {

	ctx.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
	return

}

func Example01() {

	s := web.NewServer(web.WebConf{
		Boot: "/",
	})
	//常规method模式
	s.AutoRouter(web.NormalMethod, NewTestController())
	s.Run(":10086")
}
