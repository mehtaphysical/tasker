package tasker

import (
	"fmt"
	"github.com/mehtaphysical/tasker/runner"
	"github.com/mehtaphysical/tasker/task"
	"io/ioutil"
	"time"
)

var taskHistory []*task.Task = []*task.Task{}

type Tasker struct {
	Registry          task.TaskRegistry
	Runner            runner.Runner
	TaskChan          chan *task.Task
	CompletedTaskChan chan task.CompletedTask
}

func NewTasker(taskRunner runner.Runner, workers int, tasks ...task.TaskDefinition) *Tasker {
	taskChan, completed := CreateTaskerWorkerPool(workers, taskRunner)
	return &Tasker{
		Registry:          task.NewTaskRegistry(tasks...),
		Runner:            taskRunner,
		TaskChan:          taskChan,
		CompletedTaskChan: completed,
	}
}

func (tasker *Tasker) Start(port string) {
	go StartWeb(tasker, port)
	for completedTask := range tasker.CompletedTaskChan {
		if completedTask.Error != nil {
			// handle failed tasks
			fmt.Println(fmt.Sprintf("Error executing task %s with id %s: %s", completedTask.Task.Name, completedTask.Id, completedTask.Error.Error()))
		}

		// start children tasks
		for _, t := range completedTask.Task.Children {
			taskHistory = append(taskHistory, t)
			tasker.TaskChan <- t
		}
	}
}

func (tasker *Tasker) TriggerTask(taskName string, triggerIn time.Duration) error {
	taskToTrigger, err := task.NewRootTask(taskName, tasker.Registry)
	if err != nil {
		return err
	}

	time.Sleep(triggerIn)
	taskHistory = append(taskHistory, taskToTrigger)
	tasker.TaskChan <- taskToTrigger
	return nil
}

func (tasker *Tasker) TaskHistory() []task.HistoricalTask {
	tasks := []task.HistoricalTask{}
	for _, t := range taskHistory {
		tasks = append(tasks, task.HistoricalTask{
			Id:     t.Id,
			Name:   t.Name,
			Status: t.Status,
			Output: t.Output,
		})
	}
	return tasks
}

func CreateTaskerWorkerPool(number int, taskRunner runner.Runner) (chan *task.Task, chan task.CompletedTask) {
	taskChan := make(chan *task.Task, 100)
	completed := make(chan task.CompletedTask, 100)
	for i := 0; i < number; i++ {
		go TaskerWorker(fmt.Sprintf("worker-%d", i), taskRunner, taskChan, completed)
	}
	return taskChan, completed
}

func TaskerWorker(name string, taskRunner runner.Runner, taskChan chan *task.Task, completed chan task.CompletedTask) {
	for t := range taskChan {
		for _, parent := range t.Parents {
			if parent.Status != task.Complete {
				return
			}
		}

		fmt.Println(fmt.Sprintf("Worker %s executing task %s with id %s", name, t.Name, t.Id))
		err := taskRunner.Run(t)
		output, _ := ioutil.ReadAll(t.OutputBuffer)
		t.Output += string(output)
		completed <- task.CompletedTask{
			Id:    t.Id,
			Task:  t,
			Error: err,
		}
	}
}
