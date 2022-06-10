package models

type ForumRequestDelivery struct {
	Title	string	`json:"title"`
	User	string	`json:"user"`
	Slug 	string	`json:"slug"`
}

type ForumResponse struct {
	Id      int64 	`json:"-"`
	Title	string	`json:"title"`
	User	string	`json:"user"`
	Slug 	string	`json:"slug"`
	Posts 	int		`json:"posts"`
	Threads	int 	`json:"threads"`
}