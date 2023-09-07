package forum

import (
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"time"
)

type RepoPgx struct {
	DB *pgx.ConnPool
}

func NewPgxRepository(db *pgx.ConnPool) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (r *RepoPgx) CreateForum(newForum models.ForumRequestDelivery) (forum models.ForumResponse, err error) {
	sql := `insert into "forum" ("title", "user", "slug") values ($1, $2, $3) returning "title", "user", "slug", "posts", "threads";`
	err = r.DB.QueryRow(
		sql,
		newForum.Title,
		newForum.User,
		newForum.Slug,
	).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
	)
	return
}

func (r *RepoPgx) GetForumBySlug(slug string) (forum models.ForumResponse, err error) {
	sql := `select "id", "title", "user", "slug", "posts", "threads" from "forum" where "slug" = $1;`
	err = r.DB.QueryRow(
		sql,
		slug,
	).Scan(
		&forum.Id,
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
	)
	return
}

func (r *RepoPgx) GetThreadsBySlug(slug string) (thread models.ThreadResponse, err error) {
	sql := `select "id", "title", "author", "forum", "message", "votes", "slug", "created" from "thread" where "slug" = $1;`

	err = r.DB.QueryRow(
		sql,
		slug,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	return
}

func (r *RepoPgx) GetUserByNickname(nickname string) (user models.User, err error) {
	sql := `select "nickname", "fullname", "about", "email" from "user" where "nickname" = $1;`
	err = r.DB.QueryRow(
		sql,
		nickname,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
	return
}

func (r *RepoPgx) CreateThread(newThread models.ThreadsRequest) (thread models.ThreadResponse, err error) {
	if newThread.Created.String() == "" {
		newThread.Created = time.Now()
	}

	sql := `insert into "thread" ("title", "author", "forum", "message", "slug", "created") values ($1, $2, $3, $4, $5, $6) returning "id", "title", "author", "forum", "message", "votes", "slug", "created";`
	err = r.DB.QueryRow(
		sql,
		newThread.Title,
		newThread.Author,
		newThread.Forum,
		newThread.Message,
		newThread.Slug,
		newThread.Created,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	return
}

func (r *RepoPgx) GetForumThreads(slug, limit, since, desc string) ([]models.ThreadResponse, error) {
	query := `select "id", "title", "author", "forum", "message", "votes", "slug", "created" from "thread" where "forum" = $1`
	if since != "" {
		sign := ">="
		if desc == "desc" {
			sign = "<="
		}
		query += fmt.Sprintf(` and "created" %s '%s'`, sign, since)
	}
	query += fmt.Sprintf(` order by "created" %s limit %s;`, desc, limit)

	rows, err := r.DB.Query(query, slug)
	if err != nil {
		return []models.ThreadResponse{}, err
	}
	defer rows.Close()


	threads := make([]models.ThreadResponse, 0)
	for rows.Next() {
		var thread models.ThreadResponse
		err = rows.Scan(
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
			return []models.ThreadResponse{}, err
		}

		threads = append(threads, thread)
	}
	return threads, nil
}

func (r *RepoPgx) GetUsers(forum models.ForumResponse, limit, since, desc string) ([]models.User, error) {
	users := make([]models.User, 0)

	query := `select "nickname", "about", "email", "fullname" from "user" where "id" in (select "user" from "forum_user" where forum = $1)`
	if since != "" {
		sign := ">"
		if desc == "desc" {
			sign = "<"
		}
		query += fmt.Sprintf(` and "nickname" %s '%s'`, sign, since)
	}
	query += fmt.Sprintf(` order by "nickname" %s limit %s;`, desc, limit)

	rows, err := r.DB.Query(query, forum.Id)
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Nickname,
			&user.About,
			&user.Email,
			&user.Fullname,
		)
		if err != nil {
			return []models.User{}, err
		}

		users = append(users, user)
	}
	return users, nil
}

// Method MPPR

type Data struct {
	Id		int					`db:"id"`
	Name	string				`db:"name"`
	Var1	pgtype.Int4Array	`db:"var1"`
	Var2	pgtype.Int4Array	`db:"var2"`
	Var3	pgtype.Int4Array	`db:"var3"`
}

func (r *RepoPgx) GetData() (data Data, err error) {
	err = r.DB.QueryRow(
		`select id, name, var1, var2, var3 from data where id=1`,
	).Scan(
		&data.Id,
		&data.Name,
		&data.Var1,
		&data.Var2,
		&data.Var3)
	fmt.Println(data)
	fmt.Println(err)
	return
}

func (r *RepoPgx) SetData(data DataReq) (data1 Data,err error) {
	err = r.DB.QueryRow(
		`insert into data (name, var1, var2, var3) values ($1, $2, $3, $4)
			returning id, name, var1, var2, var3;`,
		data.Name, data.Var1, data.Var2, data.Var3,
	).Scan(
		&data1.Id,
		&data1.Name,
		&data1.Var1,
		&data1.Var2,
		&data1.Var3)

	fmt.Println(data1)
	return
}