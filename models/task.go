package models

import (
	"fmt"
	"time"

	logic "github.com/FrostyDog/SAM/logic/gold_ticker"
)

type Executable func()

type Task struct {
	ticker    *time.Ticker
	fn        Executable
	closeChan chan bool
	Status    bool
}

var CurrentTask Task = createNewTask()

func (t *Task) stop() {
	t.closeChan <- true
	t.Status = false
}
func (t *Task) run() {
	t.Status = true
	go func() {
		for i := range t.ticker.C {
			t.fn()
			fmt.Println(i.Second())
			select {
			case <-t.closeChan:
				fmt.Println("stopping")
				return
			default:
				continue
			}
		}
	}()
}

func createNewTask() Task {
	var closeCh = make(chan bool)
	task := Task{ticker: time.NewTicker(1 * time.Second), fn: logic.GoldRun, Status: false, closeChan: closeCh}
	return task
}

func RunTask(t *Task) {
	t.run()
}

func StopTask(t *Task) {
	t.stop()
}
