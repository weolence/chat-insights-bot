package model

type Chat struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Name     string `json:"name"`
	Filepath string `json:"filepath"`
}
