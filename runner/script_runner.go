package runner

import (
	"fmt"
	"github.com/mehtaphysical/tasker/task"
	"os"
	"os/exec"
)

type ScriptRunner struct {
	BasePath string
	DataPath string
}

func NewScriptRunner(basePath, dataPath string) *ScriptRunner {
	if dataPath == "" {
		dataPath = basePath
	}

	return &ScriptRunner{
		BasePath: basePath,
		DataPath: dataPath,
	}
}

func (r *ScriptRunner) Run(toRun *task.Task) error {
	err := os.Mkdir(r.DataPath+"/"+toRun.Id, os.ModePerm)
	if err != nil {
		handleRunnerError(err, toRun)
	}
	cmd := exec.Command(r.DataPath + "/" + toRun.Path)
	cmd.Stdout = toRun.OutputBuffer

	cmd.Env = []string{"DATA_PATH=" + r.DataPath + "/" + toRun.Id}
	for k, v := range toRun.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	err = cmd.Start()
	toRun.Status = task.Running
	if err != nil {
		return handleRunnerError(err, toRun)
	}

	err = cmd.Wait()
	if err != nil {
		return handleRunnerError(err, toRun)
	}
	toRun.Status = task.Complete

	return nil
}
