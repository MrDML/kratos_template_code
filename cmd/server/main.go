package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"server/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	// 	flag 包实现命令行标签解析
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()
	// 初始化配置
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", log.TraceID(),
		"span_id", log.SpanID(),
	)
	// 调用 go-kratos/kratos/v2/config,创建 config 实例,并指定了来源和配置解析方法
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	// 根据配置加载资源
	if err := c.Load(); err != nil {
		panic(err)
	}

	// 将配置扫描到,通过 proto 声明的 conf struct 上
	var bc conf.Bootstrap
	/*
	type Bootstrap struct {
		Server *Server
		Data   *Data
	}
	*/
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	// 扫描后bc 就有值了

	// 通过 wire 将依赖注入,并调用 newApp 方法
	// bc.Server , bc.Data 这两个结构体都被初始化后进行创建APP
	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}

	// 清理资源
	defer cleanup()

	// start and wait for stop signal
	// 开始运行服务，阻塞等待停止信号
	if err := app.Run(); err != nil {
		panic(err)
	}
}
