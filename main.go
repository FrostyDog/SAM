package main

import (
	"github.com/FrostyDog/SAM/api"
)

func main() {
	// starts the server without starting a task. (could be started from front-end)
	api.StartServer()
}
