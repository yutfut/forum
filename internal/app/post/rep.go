package post

import "example.com/greetings/internal/app/models"

type PostRep interface {
	GetPost(id int, related []string) (postInfo models.PostInfo, err error)
	UpdatePost(id int, newPost models.UpdatePostRequest) (post models.PostResponse, err error)
}
