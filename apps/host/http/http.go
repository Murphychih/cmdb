package http

import (
	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/gin-gonic/gin"
)

// Package
// 定义hanlder

type Handler struct {
	svc host.Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Config() {
	// 从IOC中获取HostService的实例对象
	apps.GetImpl(host.AppName)
}

func (h *Handler) Name() string {
	return host.AppName
}

// 完成了 Http Handler的注册
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
	r.GET("/hosts", h.queryHost)
	r.GET("/hosts/:id", h.describeHost)
	r.PUT("/hosts/:id", h.putHost)
	r.PATCH("/hosts/:id", h.patchHost)
	r.DELETE("/hosts/:id", h.deleteHost)
}

func init() {
	apps.RegistryGin(NewHandler())
}