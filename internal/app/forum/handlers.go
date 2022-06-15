package forum

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	ForumRepo ForumRep
}


func (h *Handlers) CreateForum(ctx *fasthttp.RequestCtx) {
	var forum models.ForumRequestDelivery
	err := easyjson.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	checkForum, err := h.ForumRepo.GetForumBySlug(forum.Slug)
	if err == nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusConflict)
		body, _ := easyjson.Marshal(checkForum)
		ctx.SetBody(body)
		return
	}

	checkUser, err := h.ForumRepo.GetUserByNickname(forum.User)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	forum.User = checkUser.Nickname

	newForum, err := h.ForumRepo.CreateForum(forum)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	body, _ := easyjson.Marshal(newForum)
	ctx.SetBody(body)
}

func (h *Handlers) GetForumDetails(ctx *fasthttp.RequestCtx) {
	checkForum, err := h.ForumRepo.GetForumBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := easyjson.Marshal(checkForum)
	ctx.SetBody(body)
}

func (h *Handlers) CreateThread(ctx *fasthttp.RequestCtx) {
	checkForum, err := h.ForumRepo.GetForumBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	var thread models.ThreadsRequest
	err = easyjson.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if thread.Slug != "" {
		checkThread, err := h.ForumRepo.GetThreadsBySlug(thread.Slug)
		if err == nil {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(http.StatusConflict)
			body, _ := easyjson.Marshal(checkThread)
			ctx.SetBody(body)
			return
		}
	}

	thread.Forum = checkForum.Slug

	checkAuthor, err := h.ForumRepo.GetUserByNickname(thread.Author)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}
	thread.Author = checkAuthor.Nickname

	newThread, err := h.ForumRepo.CreateThread(thread)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	body, _ := easyjson.Marshal(newThread)
	ctx.SetBody(body)
}

func (h *Handlers) GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	_, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	limit := string(ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := string(ctx.FormValue("since"))

	desc := ""
	if string(ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	threads, err := h.ForumRepo.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := json.Marshal(threads)
	ctx.SetBody(body)
}

func (h *Handlers) ForumUsers(ctx *fasthttp.RequestCtx) {
	forum, err := h.ForumRepo.GetForumBySlug(ctx.UserValue("slug").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	limit := string(ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := string(ctx.FormValue("since"))

	desc := ""
	if string(ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	users, err := h.ForumRepo.GetUsers(forum, limit, since, desc)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := json.Marshal(users)
	ctx.SetBody(body)
}