package grpc

import (
	"context"
	"gitee.com/lshsuper/go-boot/grpc/pb"
	"testing"
	"time"
)

//TestGrpcClient 测试客户端调用
func TestGrpcClient(t *testing.T) {

	go func() {
		//启动服务端
		service := NewGrpcServer(GrpcServerConf{
			Addr: "127.0.0.1:10086",
		})

		//注册grpc业务
		pb.RegisterExampleServer(service.GetServer(), pb.NewExampleService())

		if err := service.Start(); err != nil {
			return
		}
	}()

	//等待服务端启动
	time.Sleep(time.Second * 3)

	client, err := NewGrpcClient(GrpcClientConf{
		Addr: "127.0.0.1:10086",
	})
	defer client.Close()
	if err != nil {
		t.Error(err)
		return
	}

	//调用接口
	conn := client.Conn()
	exampleService := pb.NewExampleClient(conn)

	res, err := exampleService.Say(context.Background(), &pb.SayReq{
		Message: []byte(time.Now().Format("20060102")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res.Message))

}
