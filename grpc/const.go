package grpc

import (
	"time"
)

const (
	//DialTimeout 获取grpc-conn 超时时间
	DialTimeout = 5 * time.Second

	// KeepAliveTime grpc-conn保活检测时间间隔
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout 保活检测超时时间
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize 窗口链接数
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize 窗口链接数
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize 每次发送数据包体大小限制
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize 每次接收数据包体大小限制
	MaxRecvMsgSize = 4 << 30

	MaxIdle   = 10
	MaxActive = 70
	MaxRef    = 70
	RecyTime  = time.Second * 20
	DecrTime  = time.Second * 10
)
