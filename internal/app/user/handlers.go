package user

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	UserRepo UserRep
}

func (h *Handlers) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User

	err := easyjson.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		//body, _ := easyjson.Marshal(nil)
		//ctx.SetBody(body)
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)

	user1, err1 := h.UserRepo.GetUserByNickname(user.Nickname)
	user2, err2 := h.UserRepo.GetUserByEmail(user.Email)

	if err1 == nil || err2 == nil {
		var users []models.User
		if err1 == nil {
			users = append(users, user1)
		}
		if err2 == nil && user1.About != user2.About {
			users = append(users, user2)
		}
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(users)
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	_, err = h.UserRepo.CreateUser(user)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusInternalServerError)
		//body, _ := easyjson.Marshal(nil)
		//ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	body, _ := easyjson.Marshal(user)
	ctx.SetBody(body)
}

func (h *Handlers) GetProfileByNickname(ctx *fasthttp.RequestCtx) {
	user, err := h.UserRepo.GetUserByNickname(ctx.UserValue("nickname").(string))

	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: "Can't find user by nickname:"})
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := easyjson.Marshal(user)
	ctx.SetBody(body)
}

func (h *Handlers) UpdateProfile(ctx *fasthttp.RequestCtx) {
	newUserData, err := h.UserRepo.GetUserByNickname(ctx.UserValue("nickname").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message: "Can't find user by nickname:"})
		ctx.SetBody(body)
		return
	}

	err = json.Unmarshal(ctx.PostBody(), &newUserData)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		//body, _ := easyjson.Marshal(nil)
		//ctx.SetBody(body)
		return
	}

	checkUser, err := h.UserRepo.GetUserByEmail(newUserData.Email)
	if !checkUser.IsEmpty() && checkUser.Nickname != newUserData.Nickname {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusConflict)
		//
		body, _ := easyjson.Marshal(models.MessageError{Message:"This email is already registered by user:"})
		ctx.SetBody(body)
		return
	}

	user, err := h.UserRepo.UpdateProfile(newUserData)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//body, _ := easyjson.Marshal(nil)
		//ctx.SetBody(body)
		return
	}

	if newUserData.About == "" {
		user.About = newUserData.About
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := easyjson.Marshal(user)
	ctx.SetBody(body)
}
