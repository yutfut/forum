package models

type TaskResponse struct {
	Id 					int			`json:"id"`
	Name 				string		`json:"name"`
	Description 		string		`json:"description"`
	PublicTests 		[]string	`json:"public_tests"`
	PrivateTests 		[]string	`json:"private_tests"`
	GeneratedTests 		[]string	`json:"generated_tests"`
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

type TaskTest struct {
	PrivateTests 		[]string	`json:"private_tests"`
}

type SolutionRequest struct {
	IdTask		int64			`json:"id_task"`
	Solution	string 		`json:"solution"`
}

type SourceCode struct {
	Makefile	string `json:"Makefile"`
	Main		string `json:"main.c"`
}

type SolRec struct {
	SourceCode		SourceCode `json:"sourceCode"`
	Tests 			[][]string `json:"tests"`
	BuildTimeout	int `json:"buildTimeout"`
	TestTimeout 	int `json:"testTimeout"`
}