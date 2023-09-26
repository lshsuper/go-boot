package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"
)

//grpcServer grpc 服务实例
type grpcServer struct {
	server *grpc.Server
	addr   string
}

//GrpcServerConf grpc服务端配置
type GrpcServerConf struct {
	Addr         string
	Interceptors []UnaryServerInterceptor
}

func (server *grpcServer) GetServer() *grpc.Server {
	return server.server
}

//NewGrpcServer 构建grpc服务实例
func NewGrpcServer(conf GrpcServerConf) *grpcServer {

	server := &grpcServer{
		addr: conf.Addr,
	}

	opts := []grpc.ServerOption{
		grpc.InitialWindowSize(InitialWindowSize),
		grpc.InitialConnWindowSize(InitialConnWindowSize),
		grpc.MaxSendMsgSize(MaxSendMsgSize),
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			//MinTime:             5 * time.Second,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              KeepAliveTime,
			Timeout:           KeepAliveTimeout,
			MaxConnectionIdle: 15 * time.Second,
		}), //grpc.MaxConcurrentStreams(MaxRef*MaxActive),
	}

	if conf.Interceptors != nil {
		opts = append(opts, grpc.UnaryInterceptor(WarpUnaryServerInterceptor(conf.Interceptors...)))
	}

	server.server = grpc.NewServer(opts...)

	return server
}

//Start 启动grpc 服务端
func (server *grpcServer) Start() error {

	listen, err := net.Listen("tcp", server.addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	if err = server.server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil

}
