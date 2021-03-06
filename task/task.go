package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pborman/uuid"
	"io"
)

type TaskDefinition struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Env      map[string]string `json:"env"`
	Children []string          `json:"children"`
	Parents  []string          `json:"parents"`
}

type HistoricalTask struct {
	Id     string     `json:"id"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
	Output string     `json:"output"`
}

type CompletedTask struct {
	Id    string
	Task  *Task
	Error error
}

type TaskStatus int

const (
	Pending TaskStatus = iota
	Running
	Failed
	Complete
)

func (t *TaskStatus) MarshalJSON() ([]byte, error) {
	var statusString string
	switch *t {
	case Pending:
		statusString = "pending"
	case Running:
		statusString = "running"
	case Failed:
		statusString = "failed"
	case Complete:
		statusString = "complete"
	}

	return []byte("\"" + statusString + "\""), nil
}

func (t *TaskStatus) UnmarshalJSON(p []byte) error {
	var statusString string
	err := json.Unmarshal(p, &statusString)
	if err != nil {
		return err
	}
	switch statusString {
	case "pending":
		*t = Pending
	case "running":
		*t = Running
	case "failed":
		*t = Failed
	case "complete":
		*t = Complete
	}

	return nil
}

type Task struct {
	Id           string
	Name         string
	Path         string
	Env          map[string]string
	Status       TaskStatus
	Children     []*Task
	Parents      []*Task
	OutputBuffer io.ReadWriter
	Output       string
}

type TaskRegistry map[string]TaskDefinition

func NewTaskRegistry(tasks ...TaskDefinition) TaskRegistry {
	registry := map[string]TaskDefinition{}
	for _, task := range tasks {
		registry[task.Name] = task
	}
	return TaskRegistry(registry)
}

func (r TaskRegistry) GetTaskDefinition(taskName string) TaskDefinition {
	return r[taskName]
}

func (r TaskRegistry) RegisterTask(task TaskDefinition) {
	r[task.Name] = task
}

func NewRootTask(taskName string, registry TaskRegistry) (*Task, error) {
	id := uuid.NewRandom().String()
	task, _, err := newTask(id, taskName, map[string]*Task{}, registry)
	return task, err
}

func newTask(id, taskName string, nodes map[string]*Task, registry TaskRegistry) (*Task, map[string]*Task, error) {
	taskDefinition := registry[taskName]

	parents := []*Task{}
	for _, parent := range taskDefinition.Parents {
		if parentNode, ok := nodes[parent]; !ok {
			return nil, nodes, errors.New("Unknown parent node: " + parent)
		} else {
			parents = append(parents, parentNode)
		}
	}

	task := &Task{
		Id:           id,
		Path:         taskDefinition.Path,
		Status:       Pending,
		Env:          taskDefinition.Env,
		Name:         taskDefinition.Name,
		Parents:      parents,
		OutputBuffer: bytes.NewBuffer([]byte{}),
	}

	nodes[task.Name] = task

	children := []*Task{}
	for _, child := range taskDefinition.Children {
		var childNode *Task
		var ok bool
		if childNode, ok = nodes[child]; !ok {
			var err error
			childNode, nodes, err = newTask(id, child, nodes, registry)
			if err != nil {
				return nil, nodes, err
			}
		}
		children = append(children, childNode)
	}

	task.Children = children

	return task, nodes, nil
}
