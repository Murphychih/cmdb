package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/infraboard/mcube/http/label"

	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/common/response"
	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
)


var (
	h = &handler{}
)

type handler struct {
	service host.ServiceServer
	log     *zap.Logger
}

func (h *handler) Config()  error{
	h.log = zap.L().Named(host.AppName)
	h.service = apps.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return host.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	ws.Route(ws.POST("/").To(h.CreateHost).
		Doc("create a host").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).
		Reads(host.Host{}).
		Writes(response.NewData(host.Host{})))

	ws.Route(ws.GET("/").To(h.QueryHost).
		Doc("get all hosts").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Reads(host.QueryHostRequest{}).
		Writes(response.NewData(host.HostSet{})).
		Returns(200, "OK", host.HostSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeHost).
		Doc("describe an host").
		Param(ws.PathParameter("id", "identifier of the host").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Writes(response.NewData(host.Host{})).
		Returns(200, "OK", response.NewData(host.Host{})).
		Returns(404, "Not Found", nil))
}

func init() {
	apps.RegistryRESTfulApp(h)
}
