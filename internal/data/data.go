package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"server/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client 封装数据数据客户端（初始化mysql redis）

	// 这里可以对 MySql  和  Redis 进行封装
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// 创建你们函数关闭资源
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
