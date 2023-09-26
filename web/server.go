package web

import (
	"fmt"
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"reflect"
	"strings"
)

type WebServer struct {
	*gin.Engine
	boot           string
	b              *gin.RouterGroup
	ignorePathCase bool
	traceIDKey     string
}

//WebConf 配置
type WebConf struct {
	Boot           string
	IgnorePathCase bool
	TraceIDKey     string
}

//NewWebServer web-server构造器
func NewWebServer(conf WebConf) *WebServer {
	server := &WebServer{}
	server.Engine = gin.New()
	server.b = server.Engine.Group(conf.Boot)
	server.ignorePathCase = conf.IgnorePathCase
	server.traceIDKey = conf.TraceIDKey
	return server
}

//formatCtrlName 格式化控制器名称
func (s *WebServer) formatCtrlName(ctrlName string) string {
	index := strings.Index(ctrlName, "Controller")
	if index < 0 {
		panic(fmt.Sprintf("ctrl:%s not due to rule", ctrlName))
	}
	return ctrlName[:index]
}

//ServeHTTP 拦截器
func (s *WebServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	//管理处理逻辑
	if s.ignorePathCase {
		req.URL.Path = strings.ToLower(req.URL.Path)
	}
	s.Engine.ServeHTTP(res, req)

}

func (s *WebServer) Run(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

//Use 中间件注册入口
func (s *WebServer) Use(handles ...gin.HandlerFunc) {

	s.b.Use(handles...)
}

//UseRecover 处理跨域请求,支持options访问
func (s *WebServer) UseRecover(callback func(ctx *gin.Context, err string)) *WebServer {

	s.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				msg := utils.ErrorToString(r)
				if callback != nil {
					callback(c, msg)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
				}
				//终止当前接口链路
				c.Abort()

			}
		}()
		c.Next()
	})
	return s
}

//UseTrace 处理跨域请求,支持options访问
func (s *WebServer) UseTrace(callback func(c *gin.Context, traceID string)) *WebServer {

	if len(s.traceIDKey) <= 0 {
		s.traceIDKey = defaultTraceIDKey
	}

	s.Use(func(context *gin.Context) {

		traceId := context.GetHeader(s.traceIDKey)
		if len(traceId) <= 0 {
			traceId = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
			context.Request.Header.Add(s.traceIDKey, traceId)
		}

		if callback != nil {
			callback(context, traceId)
		}
		context.Next()
	})

	return s
}

//UseCors 处理跨域请求
func (s *WebServer) UseCors() *WebServer {

	s.Use(func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin,content-Type,AccessToken,X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	})
	return s

}

//AutoRouter 自动路由适配注册
func (s *WebServer) AutoRouter(methodMode MethodMode, ctrls ...Controller) {

	//Prefix/Ignore
	for _, ctrl := range ctrls {

		curCtrl := reflect.ValueOf(ctrl)
		curCtrlObj := curCtrl.Interface().(Controller)
		prefix := curCtrlObj.Prefix()
		ctrlName := curCtrl.Type().Elem().Name()

		for i := 0; i < curCtrl.NumMethod(); i++ {

			act := curCtrl.Type().Method(i)
			actName := act.Name
			isIgnore := false
			for _, ia := range ctrl.Ignore() {
				if ia == actName {
					isIgnore = true
					break
				}
			}
			if isIgnore {
				continue
			}

			fCtrlName := s.formatCtrlName(ctrlName)
			method := methodMode.GetMethod(actName)
			router := fmt.Sprintf("/%s/%s", fCtrlName, actName)
			if len(prefix) > 0 {
				router = "/" + prefix + router
			}
			if len(s.boot) > 0 {
				router = "/" + s.boot + router
			}

			//s.printLog("method:", method, "|", "router:", router, "|", "action:", actName, "|", "controller:", ctrlNam

			ftAct := curCtrl.Method(i).Interface()
			//doc := curCtrl.Method(i).Addr().
			//	log.Println(doc)
			s.handle(method, router, func(ctx *gin.Context) {

				//构建自定义上下文
				apiCtx := ctrl.InitCtx(ctx)
				apiCtx.SetWebServer(s)

				apiCtx.SetCtrl(ctrlName)
				apiCtx.SetAct(actName)
				//全局异常拦截器
				if curCtrlObj.Recover != nil {
					defer func() {
						if e := recover(); e != nil {
							curCtrlObj.Recover(apiCtx, e)
							return
						}
					}()
				}

				//action前拦截器
				if curCtrlObj.Before != nil {
					if e := curCtrlObj.Before(apiCtx); e != nil {
						return
					}
				}

				//最终业务逻辑执行
				ctrl.CallAct(ftAct, apiCtx)

				//action后拦截器
				if curCtrlObj.After(apiCtx) != nil {
					if e := curCtrlObj.Before(apiCtx); e != nil {
						return
					}
				}

			})

		}

	}

}

//GET  g-router
func (s *WebServer) GET(router string, act func(ctx *gin.Context)) {
	s.handle(http.MethodGet, router, act)
}

//POST p-router
func (s *WebServer) POST(router string, act func(ctx *gin.Context)) {
	s.Handle(http.MethodPost, router, act)
}

//PUT pt-router
func (s *WebServer) PUT(router string, act func(ctx *gin.Context)) {
	s.handle(http.MethodPut, router, act)
}

//DELETE d-router
func (s *WebServer) DELETE(router string, act func(ctx *gin.Context)) {
	s.handle(http.MethodDelete, router, act)
}

//handle 基准路由入口
func (s *WebServer) handle(method, router string, act func(ctx *gin.Context)) {
	if s.ignorePathCase {
		router = strings.ToLower(router)
	}
	s.Handle(method, router, act)
}

//AnnotateRouter 注解路由
//@Router 路由地址
//@Method 请求方法（post/get/delete/put）
//func (s *WebServer) AnnotateRouter(ctrls ...Controller) {
//	//TODO:待实现
//
//
//	filepath.Walk("con")
//
//	f, err := parser.ParseFile(token.NewFileSet(), "", , 0)
//	if err != nil {
//		panic(err)
//	}
//
//	// 遍历AST(Scope)
//	ast.Inspect(f, func(n ast.Node) bool {
//		// 处理节点
//		astFile := n.(*ast.File)
//		for _, o := range astFile.Scope.Objects {
//			if strings.Index(o.Name, "New") == 0 {
//				//构造函数
//				funcDecl := o.Decl.(*ast.FuncDecl)
//				log.Println(funcDecl)
//
//			}
//		}
//
//		return true
//	})
//}
