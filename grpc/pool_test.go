package grpc

import (
	"context"
	"gitee.com/lshsuper/go-boot/grpc/pb"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//TestPool 测试连接池链接构建
func TestPool(t *testing.T) {
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

	pool := NewGrpcPool(GrpcPoolConf{
		Addr:      "127.0.0.1:10086",
		MaxIdle:   5,
		MaxActive: 100,
	})

	var (
		failCount, succCount int32
		wg                   = new(sync.WaitGroup)
	)

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()
			//等待通知

			client, err := pool.Get()

			if err != nil {
				atomic.AddInt32(&failCount, 1)
				return
			}
			atomic.AddInt32(&succCount, 1)
			client.Close()
		}()

	}
	//通知并发请求
	wg.Wait()

	t.Logf("succ:%d---failed:%d", succCount, failCount)

}

//TestPoolDo 测试连接池链接构建并使用
func TestPoolDo(t *testing.T) {

	go func() {
		//启动服务端
		service := NewGrpcServer(GrpcServerConf{
			Addr: "127.0.0.1:9001",
		})

		//注册grpc业务
		pb.RegisterExampleServer(service.GetServer(), pb.NewExampleService())

		if err := service.Start(); err != nil {
			return
		}
	}()

	//等待服务端启动
	time.Sleep(time.Second * 3)

	pool := NewGrpcPool(GrpcPoolConf{
		Addr:      "127.0.0.1:9001",
		MaxIdle:   5,
		MaxActive: 15,
		MaxRef:    60,
	})

	var (
		failCount, succCount int32
		wg                   = new(sync.WaitGroup)
	)

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(index int) {

			defer wg.Done()

			client, err := pool.Get()

			if err != nil {

				atomic.AddInt32(&failCount, 1)

			} else {

				//调用接口
				exampleService := pb.NewExampleClient(client.Conn())
				res, err := exampleService.Say(context.Background(), &pb.SayReq{
					Message: []byte("msg:" + strconv.Itoa(index)),
				})

				if err != nil {
					t.Log(err.Error())
				} else {
					t.Log(res)
				}
				client.Close()
				atomic.AddInt32(&succCount, 1)

			}

		}(i)

	}
	//通知并发请求
	wg.Wait()

	t.Logf("succ:%d---failed:%d---conn-len:%d", succCount, failCount, pool.connMap.len())

}
