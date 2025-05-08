package db

type Storage interface {
	AddTask(task *Task) (string, error)
	Tasks(limit int) ([]*Task, error)
	SearchTasks(search string, limit int) ([]*Task, error)
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
