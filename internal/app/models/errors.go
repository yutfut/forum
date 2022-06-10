package models

type MessageError struct {
	Message string `json:"message"`
}

const (
	USEREXISTSERROR = "user exists"
	UPDATEUSERDATAERROR = "data error"
)