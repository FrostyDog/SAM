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
	"github.com/FrostyDog/SAM/task"
	"github.com/gorilla/mux"
)

// Start /** Starts the web server listener on given host and port.
func start(host string, port int, cert string, key string) {
	r := mux.NewRouter()

	r.HandleFunc("/logs", CORS(logsHandler)).Methods("GET")
	r.HandleFunc("/status", CORS(statusHandler)).Methods("GET")
	r.HandleFunc("/status", CORS(statusChangerHandler)).Methods("Post")

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))

	if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", host, port), cert, key, r); err != nil {
		log.Fatal(err)
	}
}

func StartServer() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	cert := "/etc/letsencrypt/live/api.frostydog.space/fullchain.pem"
	key := "/etc/letsencrypt/live/api.frostydog.space/privkey.pem"
	if err != nil {
		port = 8081
	}
	start(host, port, cert, key)

}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
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
		if task.CurrentTask.Status {
			fmt.Fprintf(w, "Task already running with status: %v", task.CurrentTask.Status)
		} else {
			task.RunTask(&task.CurrentTask)
			fmt.Fprintf(w, "Task has started")
		}
	case "stopTask":
		if !task.CurrentTask.Status {
			fmt.Fprintf(w, "Task is not running. Staus: %v", task.CurrentTask.Status)
		} else {
			task.StopTask(&task.CurrentTask)
			fmt.Fprintf(w, "Task has stopped")
		}
	default:
		fmt.Fprintf(w, "Uknown action\n: %s", bodyContent.Action)
	}

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status bool `json:"status"`
	}{Status: task.CurrentTask.Status}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
