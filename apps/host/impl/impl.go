package impl

import (
	"database/sql"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/conf"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db  *sql.DB
	log *zap.Logger
	host.UnimplementedServiceServer
}

func (s *service) Config() error{
	db := conf.LoadGloabal().MySQL.GetDB()
	s.log = zap.L().Named(s.Name())
	s.db = db

	return nil
}

func (s *service) Name() string {
	return host.AppName
}

func (s *service) Registry(server *grpc.Server) {
	host.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
