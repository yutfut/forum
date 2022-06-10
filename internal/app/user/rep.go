package user

import "example.com/greetings/internal/app/models"

type UserRep interface {
	CreateUser(newUser models.User) (models.User, error)
	GetUserByNickname(nickname string) (user models.User, err error)
	UpdateProfile(user models.User) (NewUser models.User, err error)
	GetUserByEmail(email string) (user models.User, err error)
}
