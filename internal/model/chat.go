package model

type Chat struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	Name     string `json:"name"`
	Filepath string `json:"filepath"`
}
