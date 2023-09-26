package grpc

import (
	"gitee.com/lshsuper/go-boot/grpc/pb"
	"testing"
	"time"
)

//TestGrpcServer 测试服务端启动
func TestGrpcServer(t *testing.T) {

	failed := make(chan bool)
	go func() {
		//启动服务端
		service := NewGrpcServer(GrpcServerConf{
			Addr: "127.0.0.1:10086",
		})

		//注册grpc业务
		pb.RegisterExampleServer(service.GetServer(), pb.NewExampleService())

		if err := service.Start(); err != nil {
			failed <- true
			return
		}
	}()

	select {
	case _ = <-failed:
		t.Failed()
	case _ = <-time.NewTicker(time.Second * 3).C:
		t.Log("grpc server start...")
		return

	}

	return

}
