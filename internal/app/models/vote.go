package models

type VoteRequest struct {
	Nickname	string	`json:"nickname"`
	Voice 		int 	`json:"voice"`
}

type VoteResponse struct {
	Id 		int	`json:"-"`
	User	int	`json:"user"`
	Thread	int `json:"thread"`
	Voice 	int `json:"voice"`
}
