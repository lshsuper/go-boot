package simple

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

type Context struct {
	res http.ResponseWriter
	req *http.Request
	context.Context
	store          map[interface{}]interface{}
	rw             *sync.RWMutex
	actionName     string
	controllerName string
}

func (ctx *Context) ActionName() string {
	return ctx.actionName
}

func (ctx *Context) ControllerName() string {
	return ctx.controllerName
}

func (ctx *Context) Set(key interface{}, val interface{}) {

	defer ctx.rw.Unlock()
	ctx.rw.Lock()
	ctx.store[key] = val
	return

}

func (ctx *Context) Value(key interface{}) (val interface{}) {

	defer ctx.rw.RUnlock()
	ctx.rw.RLock()
	val = ctx.store[key]
	return

}

//Req 请求体
func (ctx *Context) Req() *http.Request {

	return ctx.req
}

//Res 响应体
func (ctx *Context) Res() http.ResponseWriter {
	return ctx.res
}

func (ctx *Context) Json(data interface{}) {

	ctx.res.Header().Set("Content-Type", "application/json;charset=utf-8")
	by, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, _ = ctx.res.Write(by)
	return

}

//newContext 构建上下文实例
func newContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{req: req, res: res, Context: context.Background()}
}

//GInt 获取int参数
func (ctx *Context) GInt(key string, def int) int {

	v, e := strconv.Atoi(ctx.req.URL.Query().Get(key))
	if e != nil {
		return def
	}
	return v

}

//GString 获取字符串类型
func (ctx *Context) GString(key string, def string) string {

	v := ctx.req.URL.Query().Get(key)
	if len(v) <= 0 {
		return def
	}
	return v

}

//GFloat64 获取Float64参数
func (ctx *Context) GFloat64(key string, def float64) float64 {

	v, e := strconv.ParseFloat(ctx.req.URL.Query().Get(key), 10)
	if e != nil {
		return def
	}
	return v

}

//FInt 获取int参数
func (ctx *Context) FInt(key string, def int) int {

	v, e := strconv.Atoi(ctx.req.Form.Get(key))
	if e != nil {
		return def
	}
	return v

}

//FString 获取字符串类型
func (ctx *Context) FString(key string, def string) string {

	v := ctx.req.Form.Get(key)
	if len(v) <= 0 {
		return def
	}
	return v

}

//FFloat64 获取Float64参数
func (ctx *Context) FFloat64(key string, def float64) float64 {

	v, e := strconv.ParseFloat(ctx.req.Form.Get(key), 10)
	if e != nil {
		return def
	}
	return v

}

//Proxy 代理转换  target ->来源于注册中心or配置好的地址
func (ctx *Context) Proxy(target *url.URL) {

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

	proxy.ServeHTTP(ctx.res, ctx.req)
	return

}

//func (ctx *Context) DbContext() {
//
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
//	rows, err := db.
//		rows.ColumnTypes()
//
//}
