package model

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Move string `json:"move"`
}

type MsgTask struct {
	Msg string `json:"msg"`
}
