package main

import (
	"sync"

	"github.com/FrostyDog/SAM/models"
)

func main() {
	// api.StartServer()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	models.RunTask(&models.CurrentTask)
	wg.Wait()
}
