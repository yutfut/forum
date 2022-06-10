package models

type User struct {
	Nickname	string	`json:"nickname,omitempty"`
	Fullname	string	`json:"fullname,omitempty"`
	About   	string	`json:"about,omitempty"`
	Email   	string	`json:"email,omitempty"`
}

//var UserExistsError = errors.New("exist")

func (u User) IsEmpty() bool {
	return u.Nickname == ""
}