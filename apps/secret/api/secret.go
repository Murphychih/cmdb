package api

import (
	"github.com/Murphychih/cmdb/apps/secret"
	"github.com/Murphychih/cmdb/common/response"
	"github.com/emicklei/go-restful/v3"
)

func (h *handler) CreateSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewCreateSecretRequest()
	ins, err := h.service.CreateSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
	}
	response.Success(w, ins)

}
func (h *handler) QuerySecret(r *restful.Request, w *restful.Response) {
	req := secret.NewQuerySecretRequestFromHTTP(r.Request)
	set, err := h.service.QuerySecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) DescribeSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewDescribeSecretRequest(r.PathParameter("id"))
	ins, err := h.service.DescribeSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 通过 HTTP API 对外进行数据暴力是脱敏
	ins.Data.Desense()
	response.Success(w, ins)
}

func (h *handler) DeleteSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewDeleteSecretRequestWithID(r.PathParameter("id"))
	set, err := h.service.DeleteSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}
