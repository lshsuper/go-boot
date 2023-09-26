package grpc

import (
	"context"
	"google.golang.org/grpc"
)

type UnaryClientInterceptor func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error

func WarpUnaryClientInterceptor(fns ...UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		//拦截器
		for _, fn := range fns {
			if err := fn(ctx, method, req, reply, cc, opts...); err != nil {
				return err
			}
		}

		//UnaryInvoker func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, opts ...CallOption) error
		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			return err
		}

		return nil
	}
}

type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) error

func WarpUnaryServerInterceptor(fns ...UnaryServerInterceptor) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		for _, fn := range fns {
			if err = fn(ctx, req, info); err != nil {
				return nil, err
			}
		}

		m, err := handler(ctx, req)

		return m, err
	}
}
