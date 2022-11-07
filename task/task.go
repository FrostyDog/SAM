package task

import (
	"fmt"
	"time"

	kucoin_api "github.com/FrostyDog/SAM/third-party/kucoin-api"
	"github.com/Kucoin/kucoin-go-sdk"
)

type Executable func(s *kucoin.ApiService)

type Task struct {
	ticker    *time.Ticker
	isTicker  bool
	fn        Executable
	closeChan chan bool
	Status    bool
}

func (t *Task) stop() {
	t.closeChan <- true
	t.Status = false
}

// two behavior ways based on isTicker needed or not - ws not require Ticker
func (t *Task) run() {
	// if the ticker is running
	t.Status = true
	fmt.Println("Running task")
	go func() {
		t.fn(kucoin_api.S)
		select {
		case <-t.closeChan:
			fmt.Println("stopping")
			return
		}
	}()
}

// Create a task without timer. (Timer could be implemente in task function itself)
func createNewTask(fn Executable) Task {
	var closeCh = make(chan bool)
	task := Task{isTicker: false, fn: fn, Status: false, closeChan: closeCh}
	return task
}

func RunTask(t *Task) {
	t.run()
}

func StopTask(t *Task) {
	t.stop()
}
