package impl

import (
	"database/sql"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/apps/secret"
	"github.com/Murphychih/cmdb/apps/task"
	"github.com/Murphychih/cmdb/conf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &impl{}
)

type impl struct {
	db  *sql.DB
	log *zap.Logger
	task.UnimplementedServiceServer

	secret secret.ServiceServer
	host   host.ServiceServer
}

func (s *impl) Config() error {
	db := conf.LoadGloabal().MySQL.GetDB()

	s.log = zap.L().Named(s.Name())
	s.db = db

	// 通过mock 来解耦以来 s.secret = &secretMoczap
	s.secret = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	s.host = apps.GetGrpcApp(host.AppName).(host.ServiceServer)
	
	return nil
}

func (s *impl) Name() string {
	return task.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	task.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
