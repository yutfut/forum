package post

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
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
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	related := strings.Split(fmt.Sprintf("%s", ctx.FormValue("related")), ",")

	post, err := h.PostRepo.GetPost(id, related)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find post with id:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(post)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) UpdatePost(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("id")))
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	var newPost models.UpdatePostRequest
	err = json.Unmarshal(ctx.PostBody(), &newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBody(body)
		return
	}

	var related []string
	postInfo, err := h.PostRepo.GetPost(id, related)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find post with id:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}
	oldPost := postInfo.Post
	if newPost.Message == "" || oldPost.Message == newPost.Message {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(oldPost)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}

	post, err := h.PostRepo.UpdatePost(id, newPost)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: fmt.Sprintf("Can't find post with id:")})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(post)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}