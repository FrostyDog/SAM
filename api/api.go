package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/FrostyDog/SAM/models"
	"github.com/gorilla/mux"
)

/**
Please note Start functions is a placeholder for you to start your own solution.
Feel free to drop gorilla.mux if you want and use any other solution available.
main function reads host/port from env just for an example, flavor it following your taste
*/

// Start /** Starts the web server listener on given host and port.
func start(host string, port int) {
	r := mux.NewRouter()

	r.HandleFunc("/logs", logsHandler).Methods("GET")
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/status", statusChangerHandler).Methods("Post")

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r); err != nil {
		log.Fatal(err)
	}
}
func StartServer() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	start(host, port)

}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("log.txt")
	if err != nil {
		fmt.Fprintf(w, "Hello, looks like some error is occured: %s", err)
	}
	w.Write(content)
}

func statusChangerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading body:\n%v", err)
	}

	var bodyContent models.ActionRequest
	err = json.Unmarshal(body, &bodyContent)
	if err != nil {
		fmt.Fprintf(w, "Error while parsing action item: %v", err)
	}
	w.WriteHeader(http.StatusOK)

	switch bodyContent.Action {
	case "startTask":
		if models.CurrentTask.Status {
			fmt.Fprintf(w, "Task already running with status: %v", models.CurrentTask.Status)
		} else {
			models.RunTask(&models.CurrentTask)
			fmt.Fprintf(w, "Task has started")
		}
	case "stopTask":
		if !models.CurrentTask.Status {
			fmt.Fprintf(w, "Task is not running. Staus: %v", models.CurrentTask.Status)
		} else {
			models.StopTask(&models.CurrentTask)
			fmt.Fprintf(w, "Task has stopped")
		}
	default:
		fmt.Fprintf(w, "Uknown action\n: %s", bodyContent.Action)
	}

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Current task status: %v", models.CurrentTask.Status)
}
