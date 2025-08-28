package model

import "time"

type Task struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name"`
	Move        string     `json:"move"`
	Proc        bool       `json:"proc"`
	Time        time.Time  `json:"time,omitempty"`
	TaskDueDate *time.Time `json:"duedate,omitempty"`
	UserID      int        `json:"user_id"`
}

type TaskStatus struct {
	AllTasks       int `json:"allTasks"`
	CompletedTasks int `json:"completedTasks"`
	ActiveTasks    int `json:"activeTasks"`
	OverdueTasks   int `json:"overdueTasks"`
}

type Token struct {
	Token string `json:"token"`
}

type ChekBox struct {
	Check bool `json:"check"`
}

type DueDate struct {
	Due *time.Time `json:"due"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
