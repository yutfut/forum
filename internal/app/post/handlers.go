package post

import (
	_ "encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
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
	id, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("id")))
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Error"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	related := strings.Split(fmt.Sprintf("%s", ctx.FormValue("related")), ",")
	post, err := h.PostRepo.GetPost(id, related)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Can't find post with id:"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(post)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) UpdatePost(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("id")))
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Error"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	var newPost models.UpdatePostRequest
	err = easyjson.Unmarshal(ctx.PostBody(), &newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Error"})
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBody(body)
		return
	}

	var related []string
	postInfo, err := h.PostRepo.GetPost(id, related)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Can't find post with id:"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	oldPost := postInfo.Post
	if newPost.Message == "" || oldPost.Message == newPost.Message {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(oldPost)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}

	post, err := h.PostRepo.UpdatePost(id, newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := easyjson.Marshal(models.MessageError{Message: "Can't find post with id:"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := easyjson.Marshal(post)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}