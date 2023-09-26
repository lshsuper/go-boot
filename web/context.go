package web

import (
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//Ctx 上下文接口定义
type Ctx interface {
	SetCtrl(ctrl string)
	SetAct(act string)
	GetCtrl() string
	GetAct() string
	GinCtx() *gin.Context
	SetWebServer(s *WebServer)
}

//BaseCtx 基类上下文定义
type BaseCtx struct {
	*gin.Context
	ctrl string
	act  string
	s    *WebServer
}

//NewBaseCtx 构建基类上下文
func NewBaseCtx(ctx *gin.Context) *BaseCtx {
	b := &BaseCtx{}
	b.Context = ctx
	return b
}

//SetWebServer 携带服务实例指针
func (ctx *BaseCtx) SetWebServer(s *WebServer) {
	ctx.s = s
}

//TraceID 链路trace-id(需要启用链路中间件UserTrace())
func (ctx *BaseCtx) TraceID() string {

	return ctx.GetHeader(ctx.s.traceIDKey)

}

//Browser 浏览器名称
func (ctx *BaseCtx) Browser() string {
	return utils.GetBrowser(ctx.Context)
}

//GinCtx 底层上下文
func (ctx *BaseCtx) GinCtx() *gin.Context {
	return ctx.Context
}

//SetCtrl 设置控制器名称
func (ctx *BaseCtx) SetCtrl(ctrl string) {

	ctx.ctrl = ctrl

}

//SetAct 设置action名称
func (ctx *BaseCtx) SetAct(act string) {
	ctx.act = act
}

//GetCtrl 获取ctrl控制器名称
func (ctx *BaseCtx) GetCtrl() string {
	return ctx.ctrl
}

//GetAct 获取action名称
func (ctx *BaseCtx) GetAct() string {
	return ctx.act
}

//IsAjax 是否是ajax请求
func (ctx *BaseCtx) IsAjax() bool {

	var (
		accept = ctx.GetHeader("Accept")
		xrw    = ctx.GetHeader("X-Requested-With")
	)

	if (len(accept) > 0 && strings.Contains(accept, "application/json")) ||
		(len(xrw) > 0 && strings.Contains(xrw, "XMLHttpRequest")) {
		return true
	}

	return false

}

//Proxy 代理转换  target ->来源于注册中心or配置好的地址
func (ctx *BaseCtx) Proxy(target *url.URL) {

	var (
		targetQuery = target.RawQuery

		director = func(req *http.Request) {
			req.Host = target.Host
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path, req.URL.RawPath = target.Path, target.RawPath
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		}

		proxy = &httputil.ReverseProxy{Director: director}
	)

	//代理
	proxy.ServeHTTP(ctx.Context.Writer, ctx.Context.Request)
	return

}
