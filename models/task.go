package models

import (
	"fmt"
	"sync"
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

var wg sync.WaitGroup
var CurrentTask Task = createNewTask()

func (t *Task) Stop() {
	t.closeChan <- true
	t.Status = false
}
func (t *Task) Run() {
	t.Status = true
	wg.Add(1)
	go func() {
		for i := range t.ticker.C {
			t.fn()
			fmt.Println(i.Second())
			select {
			case <-t.closeChan:
				fmt.Println("stopping")
				t.ticker.Stop()
				wg.Done()
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
	t.Run()
}

func StopTask(t *Task) {
	t.Stop()
}
