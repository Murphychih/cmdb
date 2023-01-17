package impl

import (
	"database/sql"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/resource"
	"github.com/Murphychih/cmdb/conf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// Service
	svr = &Service{}
)

type Service struct {
	db     *sql.DB
	logger *zap.Logger

	resource.UnimplementedServiceServer
}

func (s *Service) Config() error{
	db := conf.LoadGloabal().MySQL.GetDB()

	s.logger = zap.L().Named(s.Name())
	s.db = db

	return nil
}

func (s *Service) Name() string {
	return resource.AppName
}

func (s *Service) Registry(server *grpc.Server) {
	resource.RegisterServiceServer(server, svr)
}

func init() {
	apps.RegistryGrpcApp(svr)
}
