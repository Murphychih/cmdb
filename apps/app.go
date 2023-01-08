package apps

import (
	"fmt"

	"github.com/Murphychih/cmdb/apps/host"
	"github.com/gin-gonic/gin"
)

var (
	// 网络服务接口
	HostService host.Service

	// 该文件主要维护当前所有的服务的实例
	implApps = map[string]ImplService{}
	ginApps  = map[string]GinService{}
)

// IOC 容器层: 管理所有的服务的实例

// 1. HostService的实例必须注册过来, HostService才会有具体的实例, 服务启动时注册
// 2. HTTP 暴露模块, 依赖Ioc里面的HostService

type ImplService interface {
	Config()
	Name() string
}

func RegistryImpl(svc ImplService) {
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service <%s> has been registered"))
	}

	implApps[svc.Name()] = svc

	// 根据对象实现不同的网络服务接口
	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}

// 如果指定了具体类型, 就导致没增加一种类型, 多一个Get方法
// func GetHostImpl(name string) host.Service

// Get 一个Impl服务的实例：implApps
// 返回一个对象, 任何类型都可以, 使用时, 由使用方进行断言
func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}

	return nil
}

// 用户初始化 注册到Ioc容器里面的所有服务
func InitImpl() {
	for _, v := range implApps {
		v.Config()
	}
}


// 注册Gin编写的Handler
// 比如 编写了Http服务A, 只需要实现Registry方法, 就能把Handler注册给Root Router
type GinService interface {
	Config()
	Name() string
	Registry(r gin.IRouter)
}
