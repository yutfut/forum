package post

import (
	"example.com/greetings/internal/app/models"
	"github.com/jackc/pgx"
)

type RepoPgx struct {
	DB *pgx.ConnPool
}

func NewPgxRepository(db *pgx.ConnPool) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (r *RepoPgx) GetPost(id int, related []string) (postInfo models.PostInfo, err error) {
	post := models.PostResponse{}
	sql := `select "id", "parent", "author", "message", "isEdited", "forum", "thread", "created" from "post" where "id" = $1;`
	err = r.DB.QueryRow(
		sql,
		id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return
	}

	postInfo.Post = &post

	if len(related) != 0 {
		for _, q := range related {
			switch q {
			case "user":
				var user models.User
				sql := `select "nickname", "fullname", "about", "email" from "user" where "nickname" = $1;`
				err = r.DB.QueryRow(
					sql,
					post.Author,
					).Scan(
						&user.Nickname,
						&user.Fullname,
						&user.About,
						&user.Email,
					)
				if err != nil {
					return
				}
				postInfo.Author = &user
			case "forum":
				var forum models.ForumResponse
				sql := `select "title", "user", "slug", "posts", "threads" from "forum" where "slug" = $1;`
				err = r.DB.QueryRow(
					sql,
					post.Forum,
					).Scan(
						&forum.Title,
						&forum.User,
						&forum.Slug,
						&forum.Posts,
						&forum.Threads,
					)
				if err != nil {
					return
				}
				postInfo.Forum = &forum
			case "thread":
				var thread models.ThreadResponse
				sql := `select "id", "title", "author", "forum", "message", "votes", "slug", "created" from "thread" where "id" = $1;`
				err = r.DB.QueryRow(
					sql,
					post.Thread,
					).Scan(
						&thread.Id,
						&thread.Title,
						&thread.Author,
						&thread.Forum,
						&thread.Message,
						&thread.Votes,
						&thread.Slug,
						&thread.Created,
					)
				if err != nil {
					return
				}
				postInfo.Thread = &thread
			default:
				break
			}
		}
	}

	return
}

func (r *RepoPgx) UpdatePost(id int, newPost models.UpdatePostRequest) (post models.PostResponse, err error) {
	sql := `select "title", "user", "slug", "posts", "threads" from "forum" where "slug" = $1;`
	err = r.DB.QueryRow(
		sql,
		newPost.Message,
		id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	return
}