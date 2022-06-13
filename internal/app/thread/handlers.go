package thread

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	ThreadRepo ThreadRep
}

func (h *Handlers) CreatePost(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ThreadRepo.GetForumThreadBySlugOrId(slugOrId)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Tread by SlugOrId:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	var posts models.PostsRequest
	err = json.Unmarshal(ctx.PostBody(), &posts.Posts)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		body, _ := json.Marshal(err.Error())
		ctx.SetBody(body)
		return
	}

	if len(posts.Posts) == 0 {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(posts.Posts)
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(body)
		return
	}

	for _, item := range posts.Posts {
		err = h.ThreadRepo.CheckPostAuthor(item.Author)
		if err != nil {
			ctx.SetContentType("application/json")
			body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find User by nickname:")})
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(body)
			return
		}

		if item.Parent != 0 {
			err = h.ThreadRepo.CheckPostByIdAndParent(item.Parent, thread.Id)
			if err != nil {
				ctx.SetContentType("application/json")
				body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Post by Id:")})
				ctx.SetStatusCode(http.StatusConflict)
				ctx.SetBody(body)
				return
			}
		}
	}

	response, err := h.ThreadRepo.CreatePost(thread, posts)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(response.Posts)
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (h *Handlers) CreateVote(ctx *fasthttp.RequestCtx) {
	thread, err := h.ThreadRepo.GetForumThreadBySlugOrId(fmt.Sprintf("%s", ctx.UserValue("slug_or_id")))
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Tread by SlugOrId:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	var vote models.VoteRequest
	err = easyjson.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBody(body)
		return
	}

	checkUserId, err := h.ThreadRepo.GetUserIdByNickname(vote.Nickname)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find user by nickname:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	vote1, err := h.ThreadRepo.CheckUserVotes(checkUserId, thread.Id)
	if err == nil && vote.Voice == vote1.Voice {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(thread)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}
	if err != nil {
		err = h.ThreadRepo.InsertVote(checkUserId, vote, thread)
		if err != nil {
			ctx.SetContentType("application/json")
			body, _ := json.Marshal(err.Error())
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(body)
			return
		}

		thread.Votes += vote.Voice
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(thread)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
	} else {
		_, err = h.ThreadRepo.UpdateVote(vote, vote1.Id)
		if err == nil {
			thread.Votes += 2 * vote.Voice
			ctx.SetContentType("application/json")
			body, _ := easyjson.Marshal(thread)
			ctx.SetStatusCode(http.StatusOK)
			ctx.SetBody(body)
			return
		}
	}
}

func (h *Handlers) ThreadDetails(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ThreadRepo.GetForumThreadBySlugOrId(slugOrId)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Tread by SlugOrId:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(thread)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) ThreadPost(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ThreadRepo.GetForumThreadBySlugOrId(slugOrId)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Tread by SlugOrId:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	limit := fmt.Sprintf("%s", ctx.FormValue("limit"))
	if limit == "" {
		limit = "100"
	}

	since := fmt.Sprintf("%s", ctx.FormValue("since"))

	sort := fmt.Sprintf("%s", ctx.FormValue("sort"))
	if sort == "" {
		sort = "flat"
	}

	desc := ""
	if fmt.Sprintf("%s", ctx.FormValue("desc")) == "true" {
		desc = "desc"
	}

	posts, err := h.ThreadRepo.GetThreadPost(thread, limit, since, sort, desc)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(posts)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) UpdateThread(ctx *fasthttp.RequestCtx) {
	slugOrId := fmt.Sprintf("%s", ctx.UserValue("slug_or_id"))
	thread, err := h.ThreadRepo.GetForumThreadBySlugOrId(slugOrId)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find Tread by SlugOrId:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	var updateThread models.UpdateThreadsRequest
	err = easyjson.Unmarshal(ctx.PostBody(), &updateThread)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBody(body)
		return
	}

	if updateThread.Title == "" && updateThread.Message == "" {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(thread)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}
	if updateThread.Title == "" {
		updateThread.Title = thread.Title
	}
	if updateThread.Message == "" {
		updateThread.Message = thread.Message
	}

	thread, err = h.ThreadRepo.UpdateThread(thread, updateThread)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find thread by slug:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(thread)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}