package simple

type Controller interface {
	//Before Action前置操作（可以做相关逻辑：如跨域等处理）
	Before(ctx *Context) error
	//After Action后置操作（可以做相关逻辑：如跨域等处理）
	After(ctx *Context) error
	//Prefix Action接口前缀（分组名）
	Prefix() string
	//Ignore 忽略的Action
	Ignore() []string
	//Recover 异常抓取操作
	Recover(ctx *Context, err interface{})
}

//BaseController 控制器基础实现
type BaseController struct {
}

func (b *BaseController) Ignore() (lst []string) {
	return []string{"Before", "After", "Prefix", "Ignore", "Recover"}
}
func (b *BaseController) Prefix() string {
	return ""
}
func (b *BaseController) Before(ctx *Context) error {
	return nil
}
func (b *BaseController) After(ctx *Context) error {
	return nil
}
func (b *BaseController) Recover(ctx *Context, err interface{}) {

	ctx.res.WriteHeader(500)
	ctx.Json(map[string]interface{}{
		"err": "这是一个简单的错误",
		"de":  err,
	})

	return
}
