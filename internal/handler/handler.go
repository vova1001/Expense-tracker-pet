package handler

import (
	"errors"
	"sync"

	m "github.com/vova1001/Expense-tracker-pet/internal/model"
)

var (
	Task  = []m.Task{}
	NewID = 1
	mu    sync.Mutex
)

func GetTask() ([]m.Task, m.MsgTask) {
	mu.Lock()
	defer mu.Unlock()
	var msg m.MsgTask
	if len(Task) == 0 {
		msg.Msg = "no tasks"
		return nil, msg
	}
	return Task, msg
}

func PostTask(Newtask m.Task) []m.Task {
	mu.Lock()
	defer mu.Unlock()
	Newtask.ID = NewID
	NewID++
	Task = append(Task, Newtask)
	return Task
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for i, t := range Task {
		if t.ID == id {
			Task = append(Task[:i], Task[i+1:]...)
			return nil
		}
	}
	return errors.New("Task ID not found")
}

func PutTask(updated m.Task) (m.Task, error) {
	for i, tasks := range Task {
		if tasks.ID == updated.ID {
			Task[i] = updated
			return updated, nil
		}
	}
	return updated, errors.New("Task not found")
}

func PatchTask(updatedPatch map[string]interface{}, id int) (m.Task, error) {
	for i, tasks := range Task {
		if tasks.ID == id {
			name, ok := updatedPatch["name"].(string)
			if ok {
				Task[i].Name = name
			}
			move, ok := updatedPatch["move"].(string)
			if ok {
				Task[i].Move = move
			}
			return Task[i], nil
		}
	}
	return m.Task{}, errors.New("Invalid id or invalid type field")
}

func ClearAll() {
	mu.Lock()
	defer mu.Unlock()
	Task = nil
	NewID = 1
}
