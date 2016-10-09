package main

import (
	"github.com/mehtaphysical/tasker/task"
	"github.com/mehtaphysical/tasker/tasker"
)

func main() {
	defaultTasks := []task.TaskDefinition{
		{
			Name:     "ryanmehta/task1",
			Children: []string{"ryanmehta/task2"},
			Parents:  []string{},
		},
		{
			Name:     "ryanmehta/task2",
			Children: []string{},
			Parents:  []string{"ryanmehta/task1"},
		},
	}
	tasker := tasker.NewTasker(3, defaultTasks...)
	tasker.Start()
}
