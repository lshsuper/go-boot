### 数据库上下文

> 概述:基于gorm

#### 一：如何使用

##### Step1：项目内注册database(以mysql为例)

```
   database := database.Register(database.Mysql, database.DbConfig{
			Host:      host,
			Port:      port,
			UserName:  uname,
			Password:  pwd,
			DefaultDb: dbname,
		})
```

##### Step2：将第一步获取的database.IDatabase实例全局保存，通过该实例可以获取具体操作数据库的实例

```
  
  db:=database.Sessign(context.Background())
```

#### 二：database实例还具备不同能力(可自行扩展)

``` 

    //获取指定库的所有表信息
    GetTables(dbName string) []TableInfo
    //获取指定表的所有列信息 
	GetColumns(tbName string) []ColumnInfo
	
```

#### 三：关于后续，可以支持多种数据库及IDatabase的多种公共能力