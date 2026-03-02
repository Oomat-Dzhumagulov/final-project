package main

import (
	"log"
	"os"

	"github.com/Oomat-Dzhumagulov/final-project/pkg/db"
	"github.com/Oomat-Dzhumagulov/final-project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка БД: %v", err)
	}

	defer db.DB.Close()

	if err := server.Start(); err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}
}
