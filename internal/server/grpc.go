package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "server/api/helloworld/v1"
	"server/internal/conf"
	"server/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *grpc.Server {

	// 防止发生异常进行恢复
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}

	// 赋值操作 等待创建服务的时候进行调用
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	// 赋值操作 等待创建服务的时候进行调用
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	// 赋值操作 等待创建服务的时候进行调用
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	// 创建服务
	srv := grpc.NewServer(opts...)

	// 注册服务
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
