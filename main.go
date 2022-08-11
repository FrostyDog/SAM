package main

import (
	"github.com/FrostyDog/SAM/api"
	"github.com/FrostyDog/SAM/models"
)

func main() {
	// manually initialize the task and server
	models.RunTask(&models.CurrentTask)
	api.StartServer()

}
