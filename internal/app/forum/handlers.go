package forum

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type Handlers struct {
	ForumRepo ForumRep
}

func (h *Handlers) GetTask(ctx *fasthttp.RequestCtx) {
	fmt.Println("///////////////////////")
	task, err := h.ForumRepo.GetTask()
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := json.Marshal(models.MessageError{Message: err.Error()})
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := json.Marshal(task)
	ctx.SetBody(body)
}

func (h *Handlers) GetTaskById(ctx *fasthttp.RequestCtx) {
	fmt.Println("///////////////////////")

	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := json.Marshal(models.MessageError{Message: err.Error()})
		ctx.SetBody(body)
		return
	}
	//fmt.Printf(string(id))
	task, err := h.ForumRepo.GetTaskById(int64(id))
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		//
		body, _ := json.Marshal(models.MessageError{Message: err.Error()})
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	body, _ := json.Marshal(task)
	ctx.SetBody(body)
}

func (h *Handlers) SendSolution(ctx *fasthttp.RequestCtx) {
	var body models.SolutionRequest
	err := json.Unmarshal(ctx.PostBody(), &body)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	test, err := h.ForumRepo.GetTestByIdTask(body.IdTask)
	fmt.Println([]test)

	//var che = models.SourceCode{
	//	Makefile: "solution: main.o\n\tgcc -o solution main.o\nrun: solution\n\t./solution",
	//	Main: "solution: main.o\n\tgcc -o solution main.o\nrun: solution\n\t./solution"
	//}

	//a := [["1 2", "3"], ["-2 5","3"], ["10 15", "25"], ["0 0", "0"]]

	//for i := 0; i < len(test); i++ {
	//
	//var a = models.Test{
	//	Input: " ",
	//	Output: " ",
	//}
	//
	//var requ = models.SolRec {
	//	SourceCode: che,
	//	Tests: [["1 2", "3"], ["-2 5","3"], ["10 15", "25"], ["0 0", "0"]],
	//	BuildTimeout: 10,
	//	TestTimeout: 10,
	//}





	fmt.Println(test)
	fmt.Println(body)
}