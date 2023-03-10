package api

import (
	"github.com/Murphychih/cmdb/apps"
	"github.com/Murphychih/cmdb/apps/task"
	"github.com/Murphychih/cmdb/common/response"

	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/label"
	"go.uber.org/zap"
)

var (
	h = &handler{}
)

type handler struct {
	task task.ServiceServer
	log  *zap.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(task.AppName)
	h.task = apps.GetGrpcApp(task.AppName).(task.ServiceServer)
	return	nil
}

func (h *handler) Name() string {
	return task.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	ws.Route(ws.POST("/").To(h.CreatTask).
		Doc("create a task").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, "task").
		Metadata(label.Action, label.Create.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(task.CreateTaskRequst{}).
		Writes(response.NewData(task.Task{})))

	ws.Route(ws.GET("/").To(h.QueryTask).
		Doc("get all task").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, "task").
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Writes(response.NewData(task.TaskSet{})).
		Returns(200, "OK", task.TaskSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeTask).
		Doc("describe an task").
		Param(ws.PathParameter("id", "identifier of the task").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, "task").
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Writes(response.NewData(task.Task{})).
		Returns(200, "OK", response.NewData(task.Task{})).
		Returns(404, "Not Found", nil))
}

func init() {
	apps.RegistryRESTfulApp(h)
}
