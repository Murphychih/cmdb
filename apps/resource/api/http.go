package api

import (
	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/resource"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/response"
	"go.uber.org/zap"
)

var (
	h = &handler{}
)

type handler struct {
	service resource.ServiceServer
	log     *zap.Logger
}

func (h *handler) Config() {
	h.log = zap.L().Named(resource.AppName)
	h.service = apps.GetGrpcApp(resource.AppName).(resource.ServiceServer)
}

func (h *handler) Name() string {
	return resource.AppName
}

// Restful API version
// /cmdb/api/v1/resource path 包含 API Version
// 通过 Version函数定义定义
func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	// RESTful API,   resource = cmdb_resource, action: list, auth: true

	tags := []string{h.Name()}
	ws.Route(ws.GET("/search").To(h.SearchResource).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Reads(resource.SearchRequest{}).
		Writes(response.NewData(resource.SearchRequest{})).
		Returns(200, "ok", resource.ResourceSet{}))
}

func init() {
	apps.RegistryRESTfulApp(h)
}
