package {{.PkgName}}
import(
     model"{{.ModelPkg}}"
)

//{{.DaoName}}Dao 数据库操作
type  {{.DaoName}}Dao struct{

}

//New{{.DaoName}}Dao 构造器
func New{{.DaoName}}Dao()*{{.DaoName}}Dao{

     return &{{.DaoName}}Dao{}

}

//Find 根据主键id获取数据元素
func (dao*{{.DaoName}}Dao)Find(db *gorm.DB,id int)(m *model.{{.DaoName}}Model){

     m=&model.{{.DaoName}}Model{}
     db.Find(&m,"id=?",id)
     if m.ID<=0{
        m=nil
        return
     }
     return

}
