package db

import "errors"

var ErrTaskNotFound = errors.New("задача не найдена")

type Storage interface {
	AddTask(task *Task) (string, error)
	UpdateTask(task *Task) error
	Tasks(limit int) ([]*Task, error)
	SearchTasks(search string, limit int) ([]*Task, error)
	Task(id string) (*Task, error)
	Delete(id string) error
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
