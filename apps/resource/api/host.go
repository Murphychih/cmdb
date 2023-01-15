package api

import (
	"github.com/Murphychih/cmdb/apps/resource"
	"github.com/Murphychih/cmdb/common/response"
	"github.com/emicklei/go-restful/v3"
)

func (h *handler) SearchResource(r *restful.Request, w *restful.Response) {
	req, err := resource.NewSearchRequestFromHTTP(r.Request)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.Search(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
