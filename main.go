package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_pablo/pkg/api"
	"go_final_pablo/pkg/db"
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

func main() {
	if err := db.Init(getDBFile()); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%d", getPort())
	log.Printf("Сервер запущен на порту %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
