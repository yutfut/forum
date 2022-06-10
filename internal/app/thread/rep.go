package thread

import "example.com/greetings/internal/app/models"

type ThreadRep interface {
	GetForumThreadBySlug(slug string) (thread models.ThreadResponse, err error)
	GetForumThreadById(id int) (thread models.ThreadResponse, err error)
	CheckPostAuthor(author string) (err error)
	CheckPostByIdAndParent(parent, id int) (err error)
	CreatePost(thread models.ThreadResponse, newPost models.PostsRequest) (response *models.PostsResponse, err error)
	GetUserIdByNickname(nickname string) (userId int, err error)
	CheckUserVotes(user, thread int) (vote models.VoteResponse, err error)
	UpdateVote(vote models.VoteRequest, voteId int) (id int, err error)
	InsertVote(userId int, vote models.VoteRequest, thread models.ThreadResponse) (err error)
	GetThreadPost(thread models.ThreadResponse, limit, since, sort, desc string) ([]models.PostResponse, error)
	UpdateThread(oldThread models.ThreadResponse, newThread models.UpdateThreadsRequest) (thread models.ThreadResponse, err error)
}