package runner

import (
	"github.com/mehtaphysical/tasker/task"
)

type Runner interface {
	Run(toRun *task.Task) error
}

func handleRunnerError(err error, t *task.Task) error {
	t.Status = task.Failed
	return err
}
