package tasker

import (
	"encoding/json"
	"fmt"
	"github.com/mehtaphysical/tasker/task"
	"io/ioutil"
	"net/http"
	"time"
)

type triggerTask struct {
	StartInMs int    `json:"start_in_s"`
	TaskName  string `json:"task"`
}

func addRoutes(t *Tasker) {
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/tasks/register", func(w http.ResponseWriter, r *http.Request) {
		if !methodCheck(http.MethodPost, w, r) {
			return
		}

		p, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", 500)
			return
		}

		taskDefinition := task.TaskDefinition{}
		err = json.Unmarshal(p, &taskDefinition)
		if err != nil {
			http.Error(w, "Error parsing json", 500)
			return
		}

		t.Registry.RegisterTask(taskDefinition)

		jsonTask, err := json.Marshal(taskDefinition)
		if err != nil {
			http.Error(w, "Error creating response json", 500)
			return
		}
		fmt.Fprintln(w, string(jsonTask))
	})

	http.HandleFunc("/tasks/trigger", func(w http.ResponseWriter, r *http.Request) {
		if !methodCheck(http.MethodPost, w, r) {
			return
		}

		p, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", 500)
			return
		}

		trigger := triggerTask{}
		err = json.Unmarshal(p, &trigger)
		if err != nil {
			http.Error(w, "Error parsing json", 500)
			return
		}

		err = t.TriggerTask(trigger.TaskName, time.Duration(trigger.StartInMs)*time.Second)
		if err != nil {
			http.Error(w, "Error triggering task: "+err.Error(), 500)
			return
		}

		jsonTrigger, err := json.Marshal(trigger)
		if err != nil {
			http.Error(w, "Error creating response json", 500)
			return
		}
		fmt.Fprintln(w, string(jsonTrigger))
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if !methodCheck(http.MethodGet, w, r) {
			return
		}

		registry, err := json.Marshal(t.Registry)
		if err != nil {
			http.Error(w, "Error creating response json", 500)
			return
		}
		fmt.Fprintln(w, string(registry))
	})
}

func methodCheck(wanted string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != wanted {
		http.Error(w, "Must use get to register tasks", 405)
		return false
	}
	return true
}

func StartWeb(t *Tasker) {
	addRoutes(t)
	http.ListenAndServe(":8080", nil)
}
