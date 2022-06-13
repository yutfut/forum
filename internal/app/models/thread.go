package models

import "time"

type ThreadsRequest struct {
	Title	string		`json:"title"`
	Author	string		`json:"author"`
	Message	string		`json:"message"`
	Created	time.Time	`json:"created,omitempty"`
	Forum 	string		`json:"forum"`
	Slug 	string		`json:"slug"`
}

type UpdateThreadsRequest struct {
	Title	string		`json:"title"`
	Message	string		`json:"message"`
}

type ThreadResponse struct {
	Id      int			`json:"id"`
	Title   string		`json:"title"`
	Author  string   	`json:"author"`
	Forum   string   	`json:"forum"`
	Message	string		`json:"message"`
	Votes   int			`json:"votes"`
	Slug    string   	`json:"slug"`
	Created	time.Time	`json:"created"`
}
