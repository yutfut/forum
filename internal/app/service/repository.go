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

func (r *RepoPgx) GetStatus() (models.Service, error) {
	var status models.Service
	err := r.DB.QueryRow(
		`select
			(select count(*) from "forum"),
			(select count(*) from "post"),
			(select count(*) from "thread"),
			(select count(*) from "user");`,
			).Scan(
				&status.Forum,
				&status.Post,
				&status.Thread,
				&status.User,
			)
	return status, err
}

func (r *RepoPgx) Clear() error {
	err := r.DB.QueryRow(
		`truncate "user", "forum", "thread", "post", "vote", "forum_user" CASCADE;`,
		).Scan()
	return err
}