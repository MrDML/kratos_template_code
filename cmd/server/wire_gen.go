// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"server/internal/biz"
	"server/internal/conf"
	"server/internal/data"
	"server/internal/server"
	"server/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// initApp init kratos application.
// confServer 包含 HTTP Server 和 GRPC Server
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {

	// 通过配置 在 NewData 可进行初始化数据源
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}


	// 通过初始化后的数据源 初始化repository 仓库
	greeterRepo := data.NewGreeterRepo(dataData, logger)

	// 通过repository 进行封装 操作实体
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)

	// 使用操作实体 进行封装 操作实体 创建 service 相当于 Java Iservice 和 ServiceImpl
	// NewGreeterService 实现了 server端 grpc 接口方法  给客户端调用
	greeterService := service.NewGreeterService(greeterUsecase, logger)

	// 创建HTTP服务句柄 通过server 和 greeterService
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	// 创建GRC服务句柄
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	// 创建APP
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
