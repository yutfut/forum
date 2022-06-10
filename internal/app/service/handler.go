package service

import (
	"encoding/json"
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
		body, _ := json.Marshal(nil)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(status)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) ServiceClear(ctx *fasthttp.RequestCtx) {
	if err := h.ServiceRepo.Clear(); err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal("")
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
