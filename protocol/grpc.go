package protocol

import (
	"net"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/conf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPCService grpc服务
type GRPCService struct {
	svr *grpc.Server
	l   *zap.Logger
	c   *conf.Config
}

func NewGRPCService() *GRPCService {
	log := zap.L().Named("GRPC Service")

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())

	return &GRPCService{
		svr: grpcServer,
		l:   log,
		c:   conf.LoadGloabal(),
	}
}

// Start 启动GRPC服务
func (s *GRPCService) Start() {
	// 装载所有GRPC服务
	apps.LoadGrpcApp(s.svr)

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPC.GetAddr())
	if err != nil {
		s.l.Sugar().Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	s.l.Sugar().Infof("GRPC 服务监听地址: %s", s.c.App.GRPC.GetAddr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		s.l.Sugar().Errorf("start grpc service error, %s", err.Error())
		return
	}
}

// Stop 启动GRPC服务
func (s *GRPCService) Stop() error {
	s.svr.GracefulStop()
	return nil
}
