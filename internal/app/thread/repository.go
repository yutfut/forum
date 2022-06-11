package thread

import (
	"errors"
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"time"
)

type RepoPgx struct {
	DB *pgx.ConnPool
}

func NewPgxRepository(db *pgx.ConnPool) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (r *RepoPgx) GetForumThreadBySlugOrId(slug string) (thread models.ThreadResponse, err error) {
	id, _ := strconv.Atoi(slug)
	err = r.DB.QueryRow(
		`select "id", "title", "author", "forum", "message", "votes", "slug", "created" 
			from "thread" 
			where "slug" = $1 or "id" = $2;`,
		slug, id,
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
	return
}

func (r *RepoPgx) CheckPostAuthor(author string) (err error) {
	err = r.DB.QueryRow(
		`select "nickname" 
			from "user"
			where "nickname" = $1;`,
			author,
			).Scan(
				&author,
			)
	return
}

func (r *RepoPgx) CheckPostByIdAndParent(parent, id int) (err error) {
	err = r.DB.QueryRow(
		`select "id" 
			from "post" 
			where "thread" = $1 and "id" = $2;`,
			id, parent,
			).Scan(
				&id,
			)
	return
}

func (r *RepoPgx) CreatePost(thread models.ThreadResponse, newPost models.PostsRequest) (response *models.PostsResponse, err error) {
	query := `insert into "post" ("parent", "author", "message", "forum", "thread", "created") values`

	var newPostsData []interface{}
	created := time.Now()

	for i, post := range newPost.Posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", 1+i*6, 2+i*6, 3+i*6, 4+i*6, 5+i*6, 6+i*6)
		newPostsData = append(newPostsData, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, created)
	}

	query += `returning "id", "parent", "author", "message", "isEdited", "forum", "thread", "created";`

	rows, err := r.DB.Query(query, newPostsData...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	response = &models.PostsResponse{}

	for rows.Next() {
		var bufPost models.PostResponse
		err = rows.Scan(
			&bufPost.Id,
			&bufPost.Parent,
			&bufPost.Author,
			&bufPost.Message,
			&bufPost.IsEdited,
			&bufPost.Forum,
			&bufPost.Thread,
			&bufPost.Created,
		)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		response.Posts = append(response.Posts, bufPost)
	}

	return response, nil
}

func (r *RepoPgx) GetUserIdByNickname(nickname string) (userId int, err error) {
	err = r.DB.QueryRow(
		`select "id" 
			from "user"
			where "nickname" = $1;`,
			nickname,
			).Scan(
				&userId,
			)
	return
}

func (r *RepoPgx) CheckUserVotes(user, thread int) (vote models.VoteResponse, err error) {
	err = r.DB.QueryRow(
		`select "id", "user", "thread", "voice" 
			from "vote" 
			where "user" = $1 and "thread" = $2;`,
			user, thread,
			).Scan(
				&vote.Id,
				&vote.User,
				&vote.Thread,
				&vote.Voice,
			)
	return
}

func (r *RepoPgx) UpdateVote(vote models.VoteRequest, voteId int) (id int, err error) {
	err = r.DB.QueryRow(
		`update "vote" 
			set "voice" = $1 
			where "id" = $2 
			returning "id";`,
			vote.Voice, voteId,
			).Scan(
				&id,
			)
	return
}

func (r *RepoPgx) InsertVote(userId int, vote models.VoteRequest, thread models.ThreadResponse) (err error) {
	err = r.DB.QueryRow(
		`insert into "vote" ("user", "thread", "voice") 
			values ($1, $2, $3) 
			returning "user";`,
			userId, thread.Id, vote.Voice,
			).Scan(
				&userId,
			)
	return
}

func (r *RepoPgx) GetThreadPost(thread models.ThreadResponse, limit, since, sort, desc string) ([]models.PostResponse, error) {
	posts := make([]models.PostResponse, 0)

	query := `select "id", "parent", "author", "message", "isEdited", "forum", "thread", "created"
			  from "post" 
			  where "thread" = $1 `

	sign := ">"
	if desc == "desc" {
		sign = "<"
	}

	switch sort {
	case "flat":
		if since != "" {
			query += fmt.Sprintf(`and "id" %s %s `, sign, since)
		}
		query += fmt.Sprintf(`order by "created" %s, "id" %s limit %s `, desc, desc, limit)
	case "tree":
		if since != "" {
			query += fmt.Sprintf(`and "path" %s (select "path" from "post" where "id" = %s) `, sign, since)
		}
		query += fmt.Sprintf(`order by path[1] %s, path %s limit %s `, desc, desc, limit)
	case "parent_tree":
		query += `and "path" && (select array (select "id" from "post" where "thread" = $1 and "parent" = 0 `
		if since != "" {
			query += fmt.Sprintf(`and "path" %s (select path[1:1] from "post" where "id" = %s) `, sign, since)
		}
		query += fmt.Sprintf(`order by path[1] %s, path limit %s)) order by path[1] %s, path `, desc, limit, desc)
	default:
		return []models.PostResponse{}, errors.New("undefined sort type")
	}

	rows, err := r.DB.Query(query, thread.Id)
	if err != nil {
		return []models.PostResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.PostResponse
		err := rows.Scan(
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
			return []models.PostResponse{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *RepoPgx) UpdateThread(oldThread models.ThreadResponse, newThread models.UpdateThreadsRequest) (thread models.ThreadResponse, err error) {
	err = r.DB.QueryRow(
		`update "thread" 
			set "title" = $1, "message" = $2 
			where "id" = $3
			returning "id", "title", "author", "forum", "message", "votes", "slug", "created";`,
			newThread.Title, newThread.Message, oldThread.Id,
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
	return
}