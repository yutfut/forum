package models

type Service struct {
	User   	int `json:"user"`
	Forum  	int `json:"forum"`
	Thread	int `json:"thread"`
	Post  	int	`json:"post"`
}
