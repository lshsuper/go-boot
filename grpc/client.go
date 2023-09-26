package grpc

import (
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"strings"
	"time"
)

type grpcClient struct {
	conn    *grpc.ClientConn
	p       *grpcPool
	usable  chan int
	tag     string
	index   int
	lastUse time.Time
}

type GrpcClientConf struct {
	Addr              string
	UnaryInterceptors []UnaryClientInterceptor
	Pool              *grpcPool
	Index             int
	Tag               string
	MaxRef            int
}

func NewGrpcClient(conf GrpcClientConf) (*grpcClient, error) {

	return wrapGrpcClient(conf)

}

func (client *grpcClient) Tag() string {
	return client.tag
}

func (client *grpcClient) uuid() string {
	id := uuid.NewV4()
	return strings.ReplaceAll(id.String(), "-", "")

}

//wrapGrpcClient 构建一个client
func wrapGrpcClient(conf GrpcClientConf) (*grpcClient, error) {

	client := &grpcClient{p: conf.Pool, index: conf.Index}
	client.tag = client.uuid()

	conn, err := client.dial(conf)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	if conf.Pool != nil {
		client.usable = make(chan int, conf.MaxRef)
		for i := 0; i < conf.MaxRef; i++ {
			client.usable <- 1
		}

	}

	return client, nil

}

//dial 初始化客户端
func (client *grpcClient) dial(conf GrpcClientConf) (*grpc.ClientConn, error) {

	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInitialWindowSize(InitialWindowSize),
		grpc.WithInitialConnWindowSize(InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	if conf.UnaryInterceptors != nil {
		opts = append(opts, grpc.WithUnaryInterceptor(WarpUnaryClientInterceptor(conf.UnaryInterceptors...)))
	}

	conn, err := grpc.DialContext(ctx, conf.Addr, opts...)

	return conn, err

}

//Conn 获取grpc链接
func (client *grpcClient) Conn() *grpc.ClientConn {
	return client.conn
}

//Close 释放链接回池or就地关掉
func (client *grpcClient) Close() error {

	if client.p == nil {
		if client.conn != nil {
			return client.conn.Close()
		}
		return nil
	}

	select {
	case client.usable <- 1:
		return nil
	default:
		return errors.New("close is error...")
	}

}

//use 使用
func (client *grpcClient) use() error {

	select {
	case _ = <-client.usable:
		client.lastUse = time.Now()
		return nil
	default:
		return errors.New("use is error...")
	}

}
func (client *grpcClient) canRecy() bool {

	return client.lastUse.Sub(time.Now().Add(client.p.recyTime)) > 0
}
