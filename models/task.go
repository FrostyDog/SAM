package models

import (
	"fmt"
	"time"

	logic "github.com/FrostyDog/SAM/logic/grow_scraping"
	kucoin_api "github.com/FrostyDog/SAM/third-party/kucoin-api"
	"github.com/Kucoin/kucoin-go-sdk"
)

type Executable func(s *kucoin.ApiService)

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
	fmt.Println("Running task")
	go func() {
		for _ = range t.ticker.C {
			t.fn(kucoin_api.S)
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
	task := Task{ticker: time.NewTicker(10 * time.Second), fn: logic.GrowScraping, Status: false, closeChan: closeCh}
	return task
}

func RunTask(t *Task) {
	t.run()
}

func StopTask(t *Task) {
	t.stop()
}
