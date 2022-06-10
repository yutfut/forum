package forum

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	ForumRepo ForumRep
}


func (h *Handlers) CreateForum(ctx *fasthttp.RequestCtx) {
	var forum models.ForumRequestDelivery
	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		body, _ := json.Marshal(err.Error())
		ctx.SetBody(body)
		return
	}

	checkForum, err := h.ForumRepo.GetForumBySlug(forum.Slug)
	if err == nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(checkForum)
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	checkUser, err := h.ForumRepo.GetUserByNickname(forum.User)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find user by nickname: %s", forum.User)})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	forum.User = checkUser.Nickname

	newForum, err := h.ForumRepo.CreateForum(forum)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(newForum)
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (h *Handlers) GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))
	checkForum, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find forum")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(checkForum)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))

	var thread models.ThreadsRequest

	err := json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		body, _ := json.Marshal(err.Error())
		ctx.SetBody(body)
		return
	}

	if thread.Slug != "" {
		checkThread, err := h.ForumRepo.GetThreadsBySlug(thread.Slug)
		if err == nil {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(http.StatusConflict)
			body, _ := json.Marshal(checkThread)
			ctx.SetBody(body)
			return
		}
	}

	checkForum, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find forum by slug:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}
	thread.Forum = checkForum.Slug

	checkAuthor, err := h.ForumRepo.GetUserByNickname(thread.Author)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find user by nickname: %s", thread.Author)})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}
	thread.Author = checkAuthor.Nickname

	newThread, err := h.ForumRepo.CreateThread(thread)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(newThread)
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (h *Handlers) GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))
	_, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find forum by slug:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	desc := ""
	if fmt.Sprintf("%s", ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	threads, err := h.ForumRepo.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Forum by slug:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(threads)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) ForumUsers(ctx *fasthttp.RequestCtx) {
	slug := fmt.Sprintf("%s", ctx.UserValue("slug"))
	forum, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find forum:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	desc := ""
	if fmt.Sprintf("%s", ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	users, err := h.ForumRepo.GetUsers(forum, limit, since, desc)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find forum:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(users)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}