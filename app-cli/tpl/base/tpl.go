package base

import (
	"text/template"
)

type BTpl struct {
	arg Arg
}

func (b *BTpl) SArg(arg Arg) {
	b.arg = arg
}
func (b *BTpl) Print() Tpl {
	panic("请实现Print-Method")
}
func (b *BTpl) Input(inputPath string) Tpl {
	panic("请实现Input-Method")
}
func (b *BTpl) GArg() Arg {
	return b.arg
}

func (b *BTpl) FArgs() map[string]interface{} {
	panic("请实现FArgs-Method")
}

func (b *BTpl) GTplInstance(tpl string) (tmpl *template.Template, err error) {

	tmpl, err = template.New(b.arg.Tag()).Parse(tpl)
	return

}

type Arg interface {
	Tag() string
}

type Tpl interface {
	GArg() Arg
	Print() Tpl
	SArg(arg Arg)
	FArgs() map[string]interface{}
	Input(inputPath string) Tpl
}
