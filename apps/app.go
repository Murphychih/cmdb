package apps

import (
	"fmt"
	"strings"

	// "github.com/Murphychih/cmdb/apps/host"
	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	// 网络服务接口
	// HostService host.Service

	// 该文件主要维护当前所有的服务的实例
	// implApps = map[string]ImplService{}
	ginApps     = map[string]GinService{}
	grpcApps    = map[string]GRPCApp{}
	restfulApps = map[string]RESTfulApp{}
)

// IOC 容器层: 管理所有的服务的实例

// 1. HostService的实例必须注册过来, HostService才会有具体的实例, 服务启动时注册
// 2. HTTP 暴露模块, 依赖Ioc里面的HostService

// type ImplService interface {
// 	Config()
// 	Name() string
// }

// func RegistryImpl(svc ImplService) {
// 	if _, ok := implApps[svc.Name()]; ok {
// 		panic(fmt.Sprintf("service <%s> has been registered"))
// 	}

// 	implApps[svc.Name()] = svc

// 	// 根据对象实现不同的网络服务接口
// 	if v, ok := svc.(host.Service); ok {
// 		HostService = v
// 	}
// }

// 如果指定了具体类型, 就导致没增加一种类型, 多一个Get方法
// func GetHostImpl(name string) host.Service

// Get 一个Impl服务的实例：implApps
// 返回一个对象, 任何类型都可以, 使用时, 由使用方进行断言
// func GetImpl(name string) interface{} {
// 	for k, v := range implApps {
// 		if k == name {
// 			return v
// 		}
// 	}

// 	return nil
// }

// 用户初始化 注册到Ioc容器里面的所有服务
// func InitImpl() {
// 	for _, v := range implApps {
// 		v.Config()
// 	}
// }

// 注册Gin编写的Handler
// 比如 编写了Http服务A, 只需要实现Registry方法, 就能把Handler注册给Root Router
type GinService interface {
	Config()
	Name() string
	Registry(r gin.IRouter)
}

func RegistryGin(svc GinService) {
	// 服务注册到svc map当中
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}

	ginApps[svc.Name()] = svc
}

// 已经加载完成的Gin App有哪些
func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}

	return names
}

// 用户初始化 注册到Ioc容器里面的所有服务
func InitGin(r gin.IRouter) {
	// 先初始化好所有对象
	for _, v := range ginApps {
		v.Config()
	}

	// 完成Http Handler的注册
	for _, v := range ginApps {
		v.Registry(r)
	}
}

type GRPCApp interface {
	Config()
	Name() string
	Registry(*grpc.Server)
}

// RegistryService 服务实例注册
func RegistryGrpcApp(app GRPCApp) {
	// 已经注册的服务禁止再次注册
	_, ok := grpcApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("grpc app %s has registed", app.Name()))
	}

	grpcApps[app.Name()] = app
}

// LoadedGrpcApp 查询加载成功的服务
func LoadedGrpcApp() (apps []string) {
	for k := range grpcApps {
		apps = append(apps, k)
	}
	return
}

func GetGrpcApp(name string) GRPCApp {
	app, ok := grpcApps[name]
	if !ok {
		panic(fmt.Sprintf("grpc app %s not registed", name))
	}

	return app
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcApp(server *grpc.Server) {
	for _, app := range grpcApps {
		app.Registry(server)
	}
}

// HTTPService Http服务的实例
type RESTfulApp interface {
	Registry(*restful.WebService)
	Config()
	Name() string
	Version() string
}

// RegistryRESTfulApp 服务实例注册
func RegistryRESTfulApp(app RESTfulApp) {
	// 已经注册的服务禁止再次注册
	_, ok := restfulApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", app.Name()))
	}

	restfulApps[app.Name()] = app
}

// LoadedHttpApp 查询加载成功的服务
func LoadedRESTfulApp() (apps []string) {
	for k := range restfulApps {
		apps = append(apps, k)
	}
	return
}

func GetRESTfulApp(name string) RESTfulApp {
	app, ok := restfulApps[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApp 装载所有的http app
func LoadRESTfulApp(pathPrefix string, root *restful.Container) {
	for _, api := range restfulApps {
		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
		ws := new(restful.WebService)
		ws.
			Path(fmt.Sprintf("%s/%s/%s", pathPrefix, api.Version(), api.Name())).
			Consumes(restful.MIME_JSON, restful.MIME_XML).
			Produces(restful.MIME_JSON, restful.MIME_XML)

		api.Registry(ws)
		root.Add(ws)
	}
}
