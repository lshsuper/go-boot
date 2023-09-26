### GRPC工具

**10w请求同时处理可以无失败执行**

#### 一、准备工作

##### Step1：编译环境

```
https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win64.zip
将解压可执行文件配置环境变量，方便全局使用
```

##### Step2：定义proto文件

```

		syntax = "proto3";
		option go_package = "./;pb";

		package pb;

		// SayReq 请求体
		message SayReq {
		    bytes message = 1;
		}

		//SayRes 响应体
		message SayRes {
		    bytes message = 1;
		}

		//Example service 业务体
		service Example {
		   // Say is simple request.
           rpc Say(SayReq) returns (SayRes) {}
		}


```
##### Step3：编译proto文件

```
   //含义：代表编译当前目录的所有proto文件
   protoc --go_out=plugins=grpc:. *.proto

```


#### Step4：实现对应rpc方法(impl)

```
type ExampleService struct {
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) Say(context.Context, *SayReq) (*SayRes, error) {
	return &SayRes{Message: []byte("hello world 666")}, nil
}

```


### 二、服务端启动注册


```
	//启动服务端
	service := NewGrpcServer(GrpcServerConf{
		 Addr: "127.0.0.1:10086",
	})

	//注册grpc业务
	pb.RegisterExampleServer(service.GetServer(), pb.NewExampleService())

	if err := service.Start(); err != nil {
        //do something
	    return
    }

```



### 三、客户端调用（直接声明式）


```
		client, err := NewGrpcClient("127.0.0.1:10086")
		if err != nil {
             //do something
		     return
		}
        defer client.Close()   //关闭链接

		//获取连接实例+拉取service
        conn := client.Conn()
		exampleService := pb.NewExampleClient(conn)
		res, err := exampleService.Say(context.Background(), &pb.SayReq{
		     Message: []byte(time.Now().Format("20060102")),
		})
		if err != nil {
		    //do something
		    return
		}
```

### 三、客户端调用（连接池式）


```
       //声明连接池
		pool := NewGrpcPool(GrpcPoolConf{

            Addr:      "127.0.0.1:10086",
            MaxIdle:   5,
            MaxActive: 100,
	   })

		//获取连接实例+拉取service
        conn,err:= pool.Get()
        if err!=nil{
            //do something
        }
		exampleService := pb.NewExampleClient(conn.Conn())
		res, err := exampleService.Say(context.Background(), &pb.SayReq{
		     Message: []byte(time.Now().Format("20060102")),
		})
		if err != nil {
		    //do something
		    return
		}
```






