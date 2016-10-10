package runner

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/mehtaphysical/tasker/task"
	"os"
)

const (
	volumeMountPoint = "/var/taskData"
)

type Runner interface {
	Run(toRun *task.Task) error
}

type DockerRunner struct {
	Registry     string
	ImagePath    string
	DockerClient *docker.Client
}

func NewDockerRunner(dockerUrl, registry, imagePath string) (*DockerRunner, error) {
	client, err := docker.NewClient(dockerUrl)
	if err != nil {
		return nil, err
	}
	return &DockerRunner{
		Registry:     registry,
		ImagePath:    imagePath,
		DockerClient: client,
	}, nil
}

func (runner *DockerRunner) Run(toRun *task.Task) error {
	// Create a volume to store task data
	err := runner.createVolume(toRun.Id)
	if err != nil {
		return handleRunnerError(err, toRun)
	}

	envVars := []string{}
	for k, v := range toRun.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	err = runner.DockerClient.PullImage(docker.PullImageOptions{
		Repository: runner.getTaskImageName(toRun.Path),
		Tag:        "latest",
	}, docker.AuthConfiguration{})
	if err != nil {
		return handleRunnerError(err, toRun)
	}

	container, err := runner.DockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Labels: map[string]string{
				"task_id": toRun.Id,
			},
			User:  "root",
			Cmd:   []string{"./task"},
			Image: runner.getTaskImageName(toRun.Path),
			Env:   envVars,
			Volumes: map[string]struct{}{
				volumeMountPoint: struct{}{},
			},
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", toRun.Id, volumeMountPoint),
			},
			PublishAllPorts: true,
			RestartPolicy:   docker.RestartOnFailure(5),
			NetworkMode:     "bridge",
		},
	})
	if err != nil {
		return handleRunnerError(err, toRun)
	}

	err = runner.DockerClient.StartContainer(container.ID, nil)
	if err != nil {
		return handleRunnerError(err, toRun)
	}
	toRun.Status = task.Running

	err = runner.DockerClient.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		Stdout:       true,
		Logs:         true,
	})
	if err != nil {
		return handleRunnerError(err, toRun)
	}

	runner.DockerClient.WaitContainer(container.ID)
	toRun.Status = task.Complete

	runner.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:            container.ID,
		RemoveVolumes: true,
		Force:         true,
	})

	return nil
}

func (runner *DockerRunner) getTaskImageName(taskName string) string {
	imageName := ""
	if runner.Registry != "" {
		imageName += runner.Registry + "/"
	}
	if runner.ImagePath != "" {
		imageName += runner.ImagePath + "/"
	}

	imageName += taskName
	return imageName
}

func (runner *DockerRunner) createVolume(name string) error {
	volumes, err := runner.DockerClient.ListVolumes(docker.ListVolumesOptions{})
	if err != nil {
		return err
	}

	// todo: would want to fix this if we were to run a significant amount of concurrent tasks
	for _, volume := range volumes {
		if volume.Name == name {
			return nil
		}
	}

	runner.DockerClient.CreateVolume(docker.CreateVolumeOptions{
		Name: name,
	})

	return nil
}

func handleRunnerError(err error, t *task.Task) error {
	t.Status = task.Failed
	return err
}
