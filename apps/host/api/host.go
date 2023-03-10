package api

import (
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/common/pb/request"
	"github.com/Murphychih/cmdb/common/response"
	"github.com/emicklei/go-restful/v3"
)


func (h *handler) QueryHost(r *restful.Request, w *restful.Response) {
	query := host.NewQueryHostRequestFromHTTP(r.Request)
	set, err := h.service.QueryHost(r.Request.Context(), query)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) CreateHost(r *restful.Request, w *restful.Response) {
	ins := host.NewDefaultHost()
	if err := request.GetDataFromRequest(r.Request, ins); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.SyncHost(r.Request.Context(), ins)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) DescribeHost(r *restful.Request, w *restful.Response) {
	req := host.NewDescribeHostRequestWithID(r.PathParameter("id"))
	set, err := h.service.DescribeHost(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}
