package {{.PkgName}}


import(
    "gitee.com/lshsuper/go-boot/web"
)

type {{.CtrlName}}Controller struct{
     web.BaseController
}

func  New{{.CtrlName}}Controller()*{{.CtrlName}}Controller{

     return &{{.CtrlName}}Controller{}

}

//Test 测试样例
func(ctrl *{{.CtrlName}}Controller)Test(ctx *web.BaseCtx){

     //ctx.JSON(gin.H{"data":"test a little..."})

}
