package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/conf"
	"github.com/Murphychih/cmdb/protocol"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// 程序的启动以及模块的组装都在这进行

var confFile string

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "cmdb 后端API",
	Long:  "cmdb 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载程序配置
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			panic(err)
		}
		// 获取全局对象
		global := conf.LoadGloabal()
		// 初始化logger
		global.Log.LoadGloabalLogger()

		// 加载当前所有app实例
		apps.InitAllApp()

		// 获取HTTP服务管理
		manager := newManager()

		ch := make(chan os.Signal, 1)
		defer close(ch)

		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)
		go manager.WaitSignal(ch)

		return manager.start()
	},
}

// 用于管理所有需要启动的服务
// 1. HTTP服务的启动
// 2. HTTP服务的关闭
type manager struct {
	http *protocol.HttpService
	grpc *protocol.GRPCService
	l    *zap.Logger
}

func newManager() *manager {
	return &manager{
		http: protocol.NewHttpService(),
		grpc: protocol.NewGRPCService(),
		l:    zap.L().Named("CLI"),
	}
}

func (m *manager) start() error {
	m.l.Sugar().Infof("loaded grpc app: %s", apps.LoadedGrpcApp())
	m.l.Sugar().Infof("loaded http app: %s", apps.LoadedRESTfulApp())

	m.l.Sugar().Infof("loaded internal app: %s", apps.LoadedInternalApp())

	go m.grpc.Start()
	return m.http.StartHttpService()
}

// 处理来自外部的中断信号
func (m *manager) WaitSignal(ch <-chan os.Signal) {
	for v := range ch {
		switch v {
		default:
			m.l.Sugar().Infof("received signal: %v", v)
			m.http.StopHttpService()
		}
	}
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/config.toml", "cmdb API配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
