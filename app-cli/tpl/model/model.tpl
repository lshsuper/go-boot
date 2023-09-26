package {{.PName}}

//{{.MName}}Model {{.Descr}}
type {{.MName}}Model struct {
     {{range .Fields}}
     {{.Name}}  {{.DType}}  `json:"{{.JName}}" gorm:"column:{{.Column}}"`   //{{.Descr}}
     {{end}}
}

//TableName 表名
func (tb*{{.MName}}Model) TableName() string {
	return "{{.TName}}"
}
