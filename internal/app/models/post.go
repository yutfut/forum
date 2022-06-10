package models

import "time"

type PostRequest struct {
	Parent	int		`json:"parent" default0:"0"`
	Author	string	`json:"author"`
	Message	string	`json:"message"`
}

type UpdatePostRequest struct {
	Message	string	`json:"message"`
}

type PostsRequest struct {
	Posts []PostRequest
}

type PostResponse struct {
	Id       	int64     	`json:"id,omitempty"`
	Parent   	int64     	`json:"parent,omitempty"`
	Author   	string    	`json:"author"`
	Message  	string    	`json:"message"`
	IsEdited	bool      	`json:"isEdited,omitempty"`
	Forum    	string    	`json:"forum,omitempty"`
	Thread   	int32     	`json:"thread,omitempty"`
	Created		time.Time	`json:"created,omitempty"`
	Path     	int64     	`json:"-"`
}

type PostsResponse struct {
	Posts []PostResponse
}

type PostInfo struct {
	Post   	*PostResponse   `json:"post"`
	Author	*User			`json:"author,omitempty"`
	Thread 	*ThreadResponse	`json:"thread,omitempty"`
	Forum  	*ForumResponse	`json:"forum,omitempty"`
}
