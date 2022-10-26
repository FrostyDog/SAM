package task

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
	isTicker  bool
	fn        Executable
	closeChan chan bool
	Status    bool
}

var CurrentTask Task = createNewTickerTask(logic.GrowScraping)

func (t *Task) stop() {
	t.closeChan <- true
	t.Status = false
}

// two behavior ways based on isTicker needed or not - ws not require Ticker
func (t *Task) run() {
	// if the ticker is running
	if t.isTicker {
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
	} else {
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
}

// Create a task without timer. (Timer could be implemente in task function itself)
func createNewTask(fn Executable) Task {
	var closeCh = make(chan bool)
	task := Task{isTicker: false, fn: fn, Status: false, closeChan: closeCh}
	return task
}

// Creating a task with a ticker of the task level.
func createNewTickerTask(fn Executable) Task {
	var closeCh = make(chan bool)
	task := Task{isTicker: false, ticker: time.NewTicker(5 * time.Second), fn: fn, Status: false, closeChan: closeCh}
	return task
}

func RunTask(t *Task) {
	t.run()
}

func StopTask(t *Task) {
	t.stop()
}
