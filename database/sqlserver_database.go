package database

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"time"
)

type SqlServerDatabase struct {
	db  *gorm.DB
	cfg DbConfig
}

func newSqlServerDatabase(opt DbConfig) *SqlServerDatabase {
	connStr := ""
	if opt.Port <= 0 {
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=DISABLE",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.DefaultDb)
	} else {
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=DISABLE",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.Port,
			opt.DefaultDb)
	}

	//dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"

	conn, err := gorm.Open(sqlserver.Open(connStr), &gorm.Config{})
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

	return &SqlServerDatabase{db: conn, cfg: opt}
}

//Conn
func (m *SqlServerDatabase) Conn() *gorm.DB {
	return m.db
}

func (m *SqlServerDatabase) GetTables(dbName string) (tbs []TableInfo) {
	m.db.Raw(`select a.name 'table_name',b.VALUE 'table_comment' from sysobjects a
                  Left Join  sys.extended_properties b on  a.id=b.major_id and b.minor_id=0
                  where type='U'`, dbName).Scan(&tbs)
	return
}
func (m *SqlServerDatabase) GetColumns(tbName string) (columns []ColumnInfo) {
	m.db.Raw(`select 
                      b.name 'field',c.name 'type',(case   when   b.isnullable=1  then   'YES'  else   'NO'   end)is_null,
	                  (case   when   exists(SELECT  1  FROM   sysobjects   where   xtype='PK'   and   name   in   (SELECT   name   FROM   sysindexes   WHERE   indid   in(
		              SELECT   indid   FROM   sysindexkeys   WHERE   id=b.id   AND   colid=b.colid)))   then   1   else   0   end)'key',isnull(e.text,'')'default',
                  f.VALUE 'comment'
				  from sysobjects a 
				  inner join syscolumns b on a.id=b.id and a.xtype='U' and  a.name<>'dtproperties'
				  inner join systypes c on b.xusertype=c.xusertype
				  left join  syscomments   e   on   b.cdefault=e.id
				  left join  sys.extended_properties   f   on   b.id=f.major_id   and     b.colid=f.minor_id
				  where a.name=?`, tbName).Scan(&columns)
	return
}
func (m *SqlServerDatabase) Session(ctx context.Context) *gorm.DB {
	return m.db.Session(&gorm.Session{
		Context: ctx,
	})
}
