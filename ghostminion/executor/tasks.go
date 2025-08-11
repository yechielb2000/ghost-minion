package executor

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TaskType string

const (
	ScreenShotTask   TaskType = "screenshot"
	keyLoggerTask    TaskType = "keylogger"
	CommandTask      TaskType = "command"
	GetFileTask      TaskType = "getfile"
	DirListTask      TaskType = "dirlist"
	ChangeConfigTask TaskType = "changeconfig"
)

type Task struct {
	Id         string                 `json:"id"`
	Type       TaskType               `json:"type"`
	TaskParams map[string]interface{} `json:"task_params"`
}

func HandleTasks(tasks []byte) {
	for _, line := range strings.Split(strings.TrimSpace(string(tasks)), "\n") {
		if line == "" {
			continue
		}

		var task Task
		if err := json.Unmarshal([]byte(line), &task); err != nil {
			fmt.Println("Failed to parse task:", err)
			continue
		}

		switch task.Type {
		case ScreenShotTask:
			fmt.Println("Handling screenshot task:", task.Id)
			// handleScreenshot(task.TaskParams)

		case keyLoggerTask:
			fmt.Println("Handling keylogger task:", task.Id)
			// handleKeyLogger(task.TaskParams)

		case CommandTask:
			fmt.Println("Handling command task:", task.Id)
			// handleCommand(task.TaskParams)

		case GetFileTask:
			fmt.Println("Handling getfile task:", task.Id)
			// handleGetFile(task.TaskParams)

		case DirListTask:
			fmt.Println("Handling dirlist task:", task.Id)
			// handleDirList(task.TaskParams)

		case ChangeConfigTask:
			fmt.Println("Handling changeconfig task:", task.Id)
			// handleChangeConfig(task.TaskParams)
		default:
			fmt.Println("Unknown task type:", task.Type)
		}
	}
}
