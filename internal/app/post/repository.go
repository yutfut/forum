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
	var post models.PostResponse
	err = r.DB.QueryRow(
		`select "id", "parent", "author", "message", "isEdited", "forum", "thread", "created"
		from "post" 
		where "id" = $1;`,
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
				err = r.DB.QueryRow(
					`select "nickname", "fullname", "about", "email"
					from "user" 
					where "nickname" = $1;`,
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
				err = r.DB.QueryRow(
					`select "title", "user", "slug", "posts", "threads"
					from "forum" 
					where "slug" = $1;`,
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
				err = r.DB.QueryRow(
					`select "id", "title", "author", "forum", "message", "votes", "slug", "created"
					from "thread" 
					where "id" = $1;`,
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
	err = r.DB.QueryRow(
		`update "post" 
			set "message" = $1, "isEdited" = true
			where "id" = $2
			returning "id", "parent", "author", "message", "isEdited", "forum", "thread", "created";`,
			newPost.Message, id,
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