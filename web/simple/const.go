package simple

//ServerMode 服务启动模式
type ServerMode string

const (
	DEBUG   ServerMode = "debug"
	RELEASE ServerMode = "release"
)
