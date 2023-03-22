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

type TaskResponse struct {
	Id 					int			`json:"id"`
	Name 				string		`json:"name"`
	Description 		string		`json:"description"`
	PublicTests 		[]string	`json:"public_tests"`
	Difficulty  		string		`json:"difficulty"`
	CfContestId 		string		`json:"cf_contest_id"`
	CfIndex   			string		`json:"cf_index"`
	CfPoints 			string		`json:"cf_points"`
	CfRating   			string		`json:"cf_rating"`
	CfTags            	string		`json:"cf_tags"`
	TimeLimit        	string		`json:"time_limit"`
	MemoryLimitBytes 	string		`json:"memory_limit_bytes"`
	Link   				string		`json:"link"`
	TaskRu 				string		`json:"task_ru"`
	Input  				string		`json:"input"`
	Output 				string		`json:"output"`
	Note 				string		`json:"note"`
}