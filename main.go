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
			Name: "task1",
			Path: "ryanmehta/file_writer",
			Env: map[string]string{
				"TEXT": "DEFAULT_TEXT",
			},
			Children: []string{"task2"},
			Parents:  []string{},
		},
		{
			Name:     "task2",
			Path:     "ryanmehta/file_printer",
			Env:      map[string]string{},
			Children: []string{},
			Parents:  []string{"task1"},
		},
		{
			Name: "clone tasker",
			Path: "ryanmehta/git_clone",
			Env: map[string]string{
				"CLONE_URL":    "https://github.com/mehtaphysical/tasker.git",
				"CLONE_OUTPUT": "tasker",
			},
			Children: []string{"print tasker readme"},
			Parents:  []string{},
		},
		{
			Name: "print tasker readme",
			Path: "ryanmehta/file_printer",
			Env: map[string]string{
				"PRINT_FILE_PATH": "/var/taskData/tasker/README.md",
			},
			Children: []string{},
			Parents:  []string{"clone tasker"},
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
