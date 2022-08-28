package main

import (
	"github.com/FrostyDog/SAM/api"
	"github.com/FrostyDog/SAM/task"
)

func main() {
	// manually initialize the task and server
	task.RunTask(&task.CurrentTask)
	api.StartServer()

}
