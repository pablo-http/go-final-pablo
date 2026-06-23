package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540
const defaultDBFile = "scheduler.db"
const webDir = "./web"

func getPort() int {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			return p
		}
	}
	return defaultPort
}

func getDBFile() string {
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		return envFile
	}
	return defaultDBFile
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbFile := getDBFile()
	if err := initDB(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	port := getPort()
	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на порту %d, БД: %s", port, dbFile)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
