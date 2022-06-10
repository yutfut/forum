package user

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

func (r *RepoPgx) CreateUser(newUser models.User) (user models.User, err error) {
	err = r.DB.QueryRow(`INSERT INTO "user" ("nickname", "fullname", "about", "email")
									VALUES ($1, $2, $3, $4) RETURNING "nickname", "fullname", "about", "email";`,
		newUser.Nickname, newUser.Fullname, newUser.About, newUser.Email,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *RepoPgx) GetUserByNickname(nickname string) (user models.User, err error) {
	err = r.DB.QueryRow(`SELECT "nickname", "fullname", "about", "email" FROM "user"
									WHERE "nickname" = $1;`, nickname,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *RepoPgx) GetUserByEmail(email string) (user models.User, err error) {
	err = r.DB.QueryRow(`SELECT "nickname", "fullname", "about", "email" FROM "user"
									WHERE "email" = $1;`, email,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return
}

func (r *RepoPgx) UpdateProfile(user models.User) (NewUser models.User, err error) {
	err = r.DB.QueryRow(`UPDATE "user" SET "fullname" = $2, "about" = $3, "email" = $4
									WHERE "nickname" = $1 RETURNING "nickname", "fullname", "about", "email";`,
		user.Nickname, user.Fullname, user.About, user.Email,
	).Scan(&NewUser.Nickname, &NewUser.Fullname, &NewUser.About, &NewUser.Email)
	return
}


