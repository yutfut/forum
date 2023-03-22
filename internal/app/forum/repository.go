package forum

import (
	"example.com/greetings/internal/app/models"
	"github.com/jackc/pgx"
	"math/rand"
)

type RepoPgx struct {
	DB *pgx.ConnPool
}

func NewPgxRepository(db *pgx.ConnPool) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (r *RepoPgx) GetTask() (Task models.TaskResponse, err error) {
	err = r.DB.QueryRow(
		`select *
		from "task"
		where id = $1;`,
		rand.Intn(3584 - 1) + 1,
	).Scan(
		&Task.Id,
		&Task.Name,
		&Task.Description,
		&Task.PublicTests,
		&Task.PrivateTests,
		&Task.GeneratedTests,
		&Task.Difficulty,
		&Task.CfContestId,
		&Task.CfIndex,
		&Task.CfPoints,
		&Task.CfRating,
		&Task.CfTags,
		&Task.TimeLimit,
		&Task.MemoryLimitBytes,
		&Task.Link,
		&Task.TaskRu,
		&Task.Input,
		&Task.Output,
		&Task.Note,
	)
	return
}

func (r *RepoPgx) GetTaskById(IdTask int64) (Task models.TaskResponse, err error) {
	err = r.DB.QueryRow(
		`select *
		from "task"
		where id = $1;`,
		IdTask,
	).Scan(
		&Task.Id,
		&Task.Name,
		&Task.Description,
		&Task.PublicTests,
		&Task.PrivateTests,
		&Task.GeneratedTests,
		&Task.Difficulty,
		&Task.CfContestId,
		&Task.CfIndex,
		&Task.CfPoints,
		&Task.CfRating,
		&Task.CfTags,
		&Task.TimeLimit,
		&Task.MemoryLimitBytes,
		&Task.Link,
		&Task.TaskRu,
		&Task.Input,
		&Task.Output,
		&Task.Note,
	)
	return
}

func (r *RepoPgx) GetTestByIdTask(IdTask int64) (Test models.TaskTest, err error) {
	err = r.DB.QueryRow(
		`select private_tests
		from "task"
		where id = $1;`,
		IdTask,
	).Scan(
		&Test.PrivateTests,
	)
	return
}