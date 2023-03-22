package forum

import "example.com/greetings/internal/app/models"

type ForumRep interface {
	GetTask() (Task models.TaskResponse, err error)
	GetTaskById(IdTask int64) (Test models.TaskResponse, err error)
	GetTestByIdTask(IdTask int64) (Test models.TaskTest, err error)
}