package router

import(

   "gitee.com/lshsuper/go-boot/web"

)


//Register 注册路由
func Register(s *web.WebServer){

    s.AutoRouter({{.CtrlPkg}}.{{.CtrlName}}Controller)

}
