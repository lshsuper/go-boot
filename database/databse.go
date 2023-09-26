package database

import (
	"context"
	"gorm.io/gorm"
)

type IDatabase interface {
	Conn() *gorm.DB
	GetTables(dbName string) []TableInfo
	GetColumns(tbName string) []ColumnInfo
	Session(ctx context.Context) *gorm.DB
}

//Register 注册实例
func Register(dbType DbType, cfg DbConfig) IDatabase {

	switch dbType {

	case Mysql:
		return newMysqlDatabase(cfg)
	case Sqlserver:
		return newSqlServerDatabase(cfg)
	default:
		panic("不存在该数据库类型")

	}

}
