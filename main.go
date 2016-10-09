package main

import (
	"flag"
	"github.com/mehtaphysical/tasker/runner"
	"github.com/mehtaphysical/tasker/task"
	"github.com/mehtaphysical/tasker/tasker"
)

func main() {
	defaultTasks := []task.TaskDefinition{
		{
			Name: "ryanmehta/task1",
			Env: map[string]string{
				"TEXT": "DEFAULT_TEXT",
			},
			Children: []string{"ryanmehta/task2"},
			Parents:  []string{},
		},
		{
			Name:     "ryanmehta/task2",
			Env:      map[string]string{},
			Children: []string{},
			Parents:  []string{"ryanmehta/task1"},
		},
	}

	dockerUrl := flag.String("dockerUrl", "unix:///var/run/docker.sock", "connection string to connect to docker daemon")
	workers := flag.Int("workers", 3, "workers in worker pool")
	port := flag.String("port", "8080", "web server listen ports")
	flag.Parse()

	taskRunner, err := runner.NewDockerRunner(*dockerUrl, "", "")
	if err != nil {
		panic("Error initializing task runner: " + err.Error())
	}

	tasker := tasker.NewTasker(taskRunner, *workers, defaultTasks...)
	tasker.Start(*port)
}
