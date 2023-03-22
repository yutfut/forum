package forum

import (
	//"github.com/fasthttp/router"
	"github.com/georgecookeIW/fasthttprouter"
)



func SetForumRouting(r *fasthttprouter.Router, h *Handlers) {
	r.GET("/api/get_task", h.GetTask)
	r.GET("/api/get_task/details/:id", h.GetTaskById)
	r.POST("/api/send_solution", h.SendSolution)
}
