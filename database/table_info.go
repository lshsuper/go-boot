package database

type TableInfo struct {
	TableName string `gorm:"column:table_name"`
	TableComment string  `gorm:"column:table_comment"`
}
