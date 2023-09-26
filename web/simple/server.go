package simple

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type server struct {
	*http.ServeMux
	boot           string
	mode           ServerMode
	ignorePathCase bool
}

type ServerConf struct {
	Boot           string
	Mode           ServerMode
	IgnorePathCase bool
	ControllerPath string
}

func NewServer(cnf ServerConf) *server {
	return &server{
		ServeMux:       http.NewServeMux(),
		boot:           cnf.Boot,
		mode:           cnf.Mode,
		ignorePathCase: cnf.IgnorePathCase,
	}
}

type Action func(ctx *Context)

//ServeHTTP 服务
func (act Action) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	act(newContext(res, req))
}

func (s *server) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	//打印请求日志
	s.printLog("url:", req.URL.Path, "|", "method:", req.Method)
	//管理处理逻辑
	req.URL.Path = req.URL.Path + "_" + req.Method
	if s.ignorePathCase {
		req.URL.Path = strings.ToLower(req.URL.Path)
	}
	s.ServeMux.ServeHTTP(res, req)

}

func (s *server) getMethod(actName string) string {

	actName = strings.ToLower(actName)
	if strings.Index(actName, "get") >= 0 || strings.Index(actName, "query") >= 0 {
		return http.MethodGet
	}
	if strings.Index(actName, "put") >= 0 || strings.Index(actName, "edit") >= 0 {
		return http.MethodPut
	}
	if strings.Index(actName, "add") >= 0 || strings.Index(actName, "create") >= 0 {
		return http.MethodPost
	}
	if strings.Index(actName, "del") >= 0 || strings.Index(actName, "delete") >= 0 || strings.Index(actName, "remove") >= 0 {
		return http.MethodDelete
	}
	return http.MethodGet

}

func (s *server) fCtrlName(ctrlName string) string {
	index := strings.Index(ctrlName, "Controller")
	return ctrlName[:index]
}

//Register 注册控制器
func (s *server) Register(ctrls ...Controller) {
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

			fCtrlName := s.fCtrlName(ctrlName)
			method := s.getMethod(actName)
			router := fmt.Sprintf("/%s/%s", fCtrlName, actName)
			if len(prefix) > 0 {
				router = "/" + prefix + router
			}
			if len(s.boot) > 0 {
				router = "/" + s.boot + router
			}

			s.printLog("method:", method, "|", "router:", router, "|", "action:", actName, "|", "controller:", ctrlName)

			ft := curCtrl.Method(i).Interface().(func(ctx *Context))
			s.handleFunc(method, router, func(res http.ResponseWriter, req *http.Request) {

				//构建上下文
				ctx := newContext(res, req)
				ctx.actionName = actName
				ctx.controllerName = ctrlName

				//全局异常拦截器
				if curCtrlObj.Recover != nil {
					defer func() {
						if e := recover(); e != nil {
							curCtrlObj.Recover(ctx, e)
							return
						}
					}()
				}

				//action前拦截器
				if curCtrlObj.Before != nil {
					if e := curCtrlObj.Before(ctx); e != nil {
						return
					}
				}

				ft(ctx)

				//action后拦截器
				if curCtrlObj.After != nil {
					if e := curCtrlObj.After(ctx); e != nil {
						return
					}
				}

			})

		}

	}

}

////AutoRouter 自动路由
//func (s *server) AutoRouter(path string) {
//
//	filepath.Walk(path, func(fPath string, info fs.FileInfo, err error) error {
//
//		if strings.Index(info.Name(), "controller") > 0 {
//
//			f, err := parser.ParseFile(token.NewFileSet(), fPath, nil, 0)
//			if err != nil {
//				panic(err)
//			}
//
//			// 遍历AST(Scope)
//			ast.Inspect(f, func(n ast.Node) bool {
//				// 处理节点
//				astFile := n.(*ast.File)
//				for _, o := range astFile.Scope.Objects {
//					if strings.Index(o.Name, "New") == 0 {
//						//构造函数
//						funcDecl := o.Decl.(*ast.FuncDecl)
//						log.Println(funcDecl)
//
//					}
//				}
//
//				return true
//			})
//
//		}
//
//		return nil
//
//	})
//
//}

//GET get路由
func (s *server) GET(router string, act Action) {
	s.handle(http.MethodGet, router, act)
}

//POST post路由
func (s *server) POST(router string, act Action) {
	s.Handle(router+"_"+http.MethodPost, act)
}

func (s *server) PUT(router string, act Action) {
	s.handle(http.MethodPut, router, act)
}
func (s *server) DELETE(router string, act Action) {
	s.handle(http.MethodDelete, router, act)
}
func (s *server) handle(method, router string, act Action) {
	s.Handle(router+"_"+method, act)
}

func (s *server) handleFunc(method, router string, fn func(http.ResponseWriter, *http.Request)) {
	router = router + "_" + method
	if s.ignorePathCase {
		router = strings.ToLower(router)
	}
	s.HandleFunc(router, fn)
}

func (s *server) Run(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	s.printLog("lis:", addr)
	if err := server.ListenAndServe(); err != nil {
		s.printLog("lis:", addr, "启动失败", "|", "msg:", err.Error())
		return err
	}
	return nil
}

func (s *server) printLog(args ...interface{}) {
	log.Println(args...)
	return
}
