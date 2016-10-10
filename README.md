# Tasker

## Task Runner

The runner package contains a task runner interface. Other runners can be
added by implementing the interface.

### Docker Task Runner

The docker task runner runs tasks inside docker contianer by executing 
a task executable file inside the working directory. State can be saved 
and passed on to other tasks by saving files to the `/var/taskData` directory. 
This directory is a shared volume mounted to each task within a graph.

## Configuration

Tasker can be configured via command line flags:

`-dockerUrl`: Docker server endpoint (default: unix:///var/run/docker.sock)

`-scriptBasePath`: base path to use for script runner (use: $PWD/scripts to use the default scritps)

`-scriptDataPath`: path for data sharing between containers (defaults to scriptBasePath)

`-workers`: Number of workers in the worker pool (default: 3)

`-port`: Web server listen port (default: 8080)

## Running tasker

### Mac

To run tasker on osx use the darwin binary passing along any configuration.
For example, to use the script task runner do:

`./tasker_darwin_amd64 -scriptBasePath=$PWD/scripts`

## Todos

* import a graph library to check for cycles
* Better error checking and handling
* Improve web ui
* Add triggers package and add more task triggers (cron, watch, etc.)
