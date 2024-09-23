package todo

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HxX2/todo/pkg/pprint"
)

type Task struct {
	Name string
	Done bool
}

type Todo struct {
	filePath string
	Tasks    []Task

	ListDone     bool
	ListUndone   bool
	ShowProgress bool
}

// Init initializes the todo struct and creates the todo.txt file if it doesn't exist
func Init() *Todo {
	todo := new(Todo)

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "todo")
	filePath := filepath.Join(configDir, "todo.txt")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			pprint.Error(fmt.Sprintf("Can't create config directory\n%s\n", err))
			return nil
		}

		file, err := os.Create(filePath)
		defer file.Close()
		if err != nil {
			pprint.Error(fmt.Sprintf("Can't create todo.txt file\n%s\n", err))
			return nil
		}

	} else if err != nil {
		pprint.Error(fmt.Sprintf("Can't check file\n%s\n", err))
		return nil
	}

	todo.filePath = filePath
	todo.ListDone = true
	todo.ListUndone = true
	todo.ShowProgress = true

	// Load tasks from file
	todo.loadTasks()

	return todo
}

// loadTasks reads tasks from the file and populates the Todo struct
func (t *Todo) loadTasks() {
	file, err := os.Open(t.filePath)
	if err != nil {
		pprint.Error(fmt.Sprintf("Can't read todo.txt file\n%s\n", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			// Parse task
			done := false
			if strings.HasPrefix(line, "[x] ") {
				done = true
				line = strings.TrimPrefix(line, "[x] ")
			} else {
				line = strings.TrimPrefix(line, "[ ] ")
			}
			t.Tasks = append(t.Tasks, Task{Name: line, Done: done})
		}
	}
}

// saveTasks writes the current tasks to the file
func (t *Todo) saveTasks() error {
	file, err := os.Create(t.filePath)
	if err != nil {
		return fmt.Errorf("can't save tasks to file: %v", err)
	}
	defer file.Close()

	for _, task := range t.Tasks {
		status := "[ ]"
		if task.Done {
			status = "[x]"
		}
		_, err := file.WriteString(fmt.Sprintf("%s %s\n", status, task.Name))
		if err != nil {
			return fmt.Errorf("can't write task to file: %v", err)
		}
	}
	return nil
}

// AddTask adds a new task to the list
func (t *Todo) AddTask(taskName string) error {
	t.Tasks = append(t.Tasks, Task{Name: taskName, Done: false})
	return t.saveTasks()
}

// RemTask removes a task by index
func (t *Todo) RemTask(taskNum int) error {
	if taskNum < 1 || taskNum > len(t.Tasks) {
		return errors.New("task number out of range")
	}
	t.Tasks = append(t.Tasks[:taskNum-1], t.Tasks[taskNum:]...)
	return t.saveTasks()
}

// ToggleTask toggles the "done" status of a task by index
func (t *Todo) ToggleTask(taskNum int) error {
	if taskNum < 1 || taskNum > len(t.Tasks) {
		return errors.New("task number out of range")
	}
	t.Tasks[taskNum-1].Done = !t.Tasks[taskNum-1].Done
	return t.saveTasks()
}
