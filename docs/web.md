### Web mvc框架

**基于gin做二次封装并抽象化请求上下文、控制器，自动路由注册等，方便扩展使用**


#### 一、简单使用


##### Step1：基础定义

```
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

```

##### Step2：注册启动服务

```


    server:=web.NewServer(web.WebConf{
           Boot:"/",
    })
    //注册控制器
    server.AutoRouter(web.NormalMethod, NewTestController())
    server.Run(":10086")



```


#### 二、自定义上下文


##### Step1：基础定义

```
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

```

##### Step2：样例使用

```
    type Example02Controller struct {
        ApiController
    }
    
    func NewExample02Controller() *Example02Controller {
        return &Example02Controller{}
    }
    
    func (ctrl *Example02Controller) Test(ctx *ApiCtx) {
    
        ctx.JSON(http.StatusOK, gin.H{
            "data": "ok",
        })
        return
    
    }
```

##### Step3：启动

```


    s := web.NewServer(web.WebConf{
		Boot: "/",
	})
	//常规method模式
	s.AutoRouter(web.NormalMethod, NewExample02Controller())
	s.Run(":10086")



```

