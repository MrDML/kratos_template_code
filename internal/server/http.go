package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "server/api/helloworld/v1"
	"server/internal/conf"
	"server/internal/service"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {

	// 切片初始化
	var opts = []http.ServerOption{
		// 中间件 Middleware is HTTP/gRPC transport middleware.  是 HTTP / gRPC 进行转换使用
		http.Middleware(
			// 恢复是一个服务器中间件，可以从任何恐慌中恢复。
			recovery.Recovery(),
		),
	}
	// // 给服务赋值Network操作
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	// 给服务赋值Addr操作
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	// 给服务赋值Timeout操作
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 通过 opts 创建服务
	srv := http.NewServer(opts...)

	// 注册服务
	v1.RegisterGreeterHTTPServer(srv, greeter)


	return srv
}
