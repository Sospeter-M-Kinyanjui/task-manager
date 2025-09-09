package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"user_id"`
}
