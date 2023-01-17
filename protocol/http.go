package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/conf"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 该模块旨在优雅的启动服务以及关闭服务

type HttpService struct {
	server *http.Server
	l      *zap.Logger
	r      gin.IRouter
}

func NewHttpService() *HttpService {
	r := gin.Default()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.LoadGloabal().App.HTTP.GetAddr(),
		Handler:           r,
	}

	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

func (s *HttpService) StartHttpService() error {
	// 加载handler
	apps.InitGin(s.r)

	// 已加载App的日志信息
	apps := apps.LoadedGinApps()
	s.l.Sugar().Infof("loaded gin apps: %s", apps)

	// 启动服务
	// 若服务正常关闭则打印日志
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Sugar().Infof("service stopped successfully")
			return nil
		}
		return fmt.Errorf("start service failed: %v", err)
	}

	return nil
}

func (s *HttpService) StopHttpService() {
	s.l.Info("start to shut down HttpService")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Sugar().DPanicf("shut down HttpService error: %v", err)
	}
}
