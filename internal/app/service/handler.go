package service

import (
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	ServiceRepo ServiceRep
}

func (h *Handlers) ServiceStatus(ctx *fasthttp.RequestCtx) {
	status, err := h.ServiceRepo.GetStatus()
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(nil)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(status)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) ServiceClear(ctx *fasthttp.RequestCtx) {
	if err := h.ServiceRepo.Clear(); err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(nil)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
