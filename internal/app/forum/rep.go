package forum

import "example.com/greetings/internal/app/models"

type ForumRep interface {
	CreateForum(newForum models.ForumRequestDelivery) (models.ForumResponse, error)
	GetForumBySlug(slug string) (forum models.ForumResponse, err error)
	GetUserByNickname(nickname string) (user models.User, err error)
	CreateThread(newThread models.ThreadsRequest) (thread models.ThreadResponse, err error)
	GetForumThreads(slug, limit, since, desc string) ([]models.ThreadResponse, error)
	GetThreadsBySlug(slug string) (thread models.ThreadResponse, err error)
	GetUsers(forum models.ForumResponse, limit, since, desc string) ([]models.User, error)
	GetTask() (forum models.TaskResponse, err error)
}