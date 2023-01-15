package impl

import (
	"database/sql"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/apps/secret"
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

	host host.ServiceServer
	secret.UnimplementedServiceServer
}

func (s *impl) Config()  {
	db := conf.LoadGloabal().MySQL.GetDB()

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.host = apps.GetGrpcApp(host.AppName).(host.ServiceServer)
}

func (s *impl) Name() string {
	return secret.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	secret.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
