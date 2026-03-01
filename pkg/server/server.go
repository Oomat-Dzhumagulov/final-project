package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Oomat-Dzhumagulov/final-project/pkg/api"
)

func Start() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	mux := http.NewServeMux()

	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux)

	fmt.Printf("Сервер запущен на порте: %s\n", port)

	return http.ListenAndServe(":"+port, mux)
}
