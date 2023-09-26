package database

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type MysqlDatabase struct {
	db  *gorm.DB
	cfg DbConfig
}

type DbConfig struct {
	Host      string
	Port      int
	UserName  string
	Password  string
	DefaultDb string
}

func newMysqlDatabase(opt DbConfig) *MysqlDatabase {

	connStr := ""
	if opt.Port <= 0 {
		connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.DefaultDb)
	} else {
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.Port,
			opt.DefaultDb)
	}
	conn, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("数据库初始异常")
	}

	db, _ := conn.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.SetConnMaxLifetime(time.Hour)

	return &MysqlDatabase{db: conn, cfg: opt}
}

//Conn
func (m *MysqlDatabase) Conn() *gorm.DB {
	return m.db
}

func (m *MysqlDatabase) GetTables(dbName string) (tbs []TableInfo) {
	m.db.Raw(`select table_name,table_comment from information_schema.tables where table_schema=?`, dbName).Scan(&tbs)
	return
}
func (m *MysqlDatabase) GetColumns(tbName string) (columns []ColumnInfo) {
	m.db.Raw("SELECT COLUMN_NAME 'field',COLUMN_TYPE 'type',COLUMN_KEY 'key',IS_NULLABLE 'is_null',COLUMN_COMMENT 'comment' FROM information_schema. COLUMNS WHERE  TABLE_NAME =?", tbName).Scan(&columns)
	return
}

func (m *MysqlDatabase) Session(ctx context.Context) *gorm.DB {

	return m.db.Session(&gorm.Session{
		Context: ctx,
	})

}
