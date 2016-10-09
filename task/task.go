package task

import (
	"errors"
	"github.com/pborman/uuid"
)

type TaskDefinition struct {
	Name     string            `json:"name"`
	Env      map[string]string `json:"env"`
	Children []string          `json:"children"`
	Parents  []string          `json:"parents"`
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

type Task struct {
	Id       string
	Name     string
	Env      map[string]string
	Status   TaskStatus
	Children []*Task
	Parents  []*Task
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
		Id:      id,
		Status:  Pending,
		Env:     taskDefinition.Env,
		Name:    taskDefinition.Name,
		Parents: parents,
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
