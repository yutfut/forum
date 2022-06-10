package user

import "github.com/fasthttp/router"

func SetServiceRouting(r *router.Router, h *Handlers) {
	r.POST("/api/user/{nickname}/create", h.CreateUser)
	r.GET("/api/user/{nickname}/profile", h.GetProfileByNickname)
	r.POST("/api/user/{nickname}/profile", h.UpdateProfile)
}