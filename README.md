# Tasker

## Task Runner

The runner package contains a task runner interface. Other runners can be
added by implementing the interface.

### Docker Task Runner

The docker task runner runs tasks inside docker contianer by executing 
a task.sh file inside the working directory. State can be saved and passed
on to other tasks by saving files to the `/var/taskData` directory. This
directory is a shared volume mounted to each task within a graph.

## Configuration

Tasker can be configured via command line flags:

`-dockerUrl`: Docker server endpoint (default: unix:///var/run/docker.sock)

`workers`: Number of workers in the worker pool (default: 3)

`port`: Web server listen port (default: 8080)

## Todos

* import a graph library to check for cycles
* Better error checking and handling
* Improve web ui
* Add more task runners (e.g. a local task runner that executes local scripts)
* Add triggers package and add more task triggers (cron, watch, etc.)
