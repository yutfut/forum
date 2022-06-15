package post

import (
	_ "encoding/json"
	"example.com/greetings/internal/app/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	PostRepo PostRep
}

func (h *Handlers) PostDetails(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	related := strings.Split(string(ctx.FormValue("related")), ",")
	post, err := h.PostRepo.GetPost(id, related)
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
	body, _ := easyjson.Marshal(post)
	ctx.SetBody(body)
}

func (h *Handlers) UpdatePost(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	var newPost models.UpdatePostRequest
	err = easyjson.Unmarshal(ctx.PostBody(), &newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	var related []string
	postInfo, err := h.PostRepo.GetPost(id, related)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		ctx.SetBody(body)
		return
	}

	oldPost := postInfo.Post
	if newPost.Message == "" || oldPost.Message == newPost.Message {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		body, _ := easyjson.Marshal(oldPost)
		ctx.SetBody(body)
		return
	}

	post, err := h.PostRepo.UpdatePost(id, newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(models.MessageError{Message: ""})
		//ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := easyjson.Marshal(post)
	ctx.SetBody(body)
	return
}