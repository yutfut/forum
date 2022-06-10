package thread

import "github.com/fasthttp/router"

func SetThreadRouting(r *router.Router, h *Handlers) {
	r.POST("/api/thread/{slug_or_id}/create", h.CreatePost)
	r.POST("/api/thread/{slug_or_id}/vote", h.CreateVote)
	r.GET("/api/thread/{slug_or_id}/details", h.ThreadDetails)
	r.GET("/api/thread/{slug_or_id}/posts", h.ThreadPost)
	r.POST("/api/thread/{slug_or_id}/details", h.UpdateThread)
}

