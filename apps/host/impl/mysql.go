package impl

import (
	"database/sql"

	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/conf"
	"go.uber.org/zap"
)

// 该服务的注册将在IOC层 apps/all/impl.go 中完成
// 如mysql 的驱动加载的实现方式
// sql 这个库, 是一个框架, 驱动是 引入依赖的时候加载的
// 把 app模块，比作一个驱动, ioc比作框架
// _ import app， 该app就注册到ioc层
// IOC的实现在 apps/app.go

type HostService struct {
	l  *zap.Logger
	db *sql.DB
}

func NewHostSerivice() *HostService {
	return &HostService{
		l:  zap.L().Named("Host"),
		db: conf.C().MySQL.GetDB(),
	}
}

func (i *HostService) Config() {
	// 初始化HostService
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.GetDB()
}

func (i *HostService) Name() string {
	return host.AppName
}

// _ import app 自动执行注册逻辑
func init() {
	//  对象注册到ioc层
	apps.RegistryImpl(NewHostSerivice())
}
