package service

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

func (r *RepoPgx) GetStatus() (status models.Service, err error) {
	sql := `select (select count(*) from "forum"), (select count(*) from "post"), (select count(*) from "thread"), (select count(*) from "user");`
	err = r.DB.QueryRow(
		sql,
	).Scan(
		&status.Forum,
		&status.Post,
		&status.Thread,
		&status.User,
	)
	return
}

func (r *RepoPgx) Clear() (err error) {
	sql := `truncate "user", "forum", "thread", "post", "vote", "forum_user" CASCADE;`
	err = r.DB.QueryRow(
		sql,
		).Scan()
	return
}