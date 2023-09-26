package database

type DbType int

const (
	_DbType = iota
	Mysql
	Sqlserver
)

func (e DbType) GetRemark() string {
	switch e {
	case Mysql:
		return "mysql"
	case Sqlserver:
		return "sqlserver"
	default:
		return ""
	}
}
