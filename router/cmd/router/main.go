package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sumanhara9054/v1/dynamo-lite/pkg/api"
	
	
)

func main() {

	
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hii hallow"))
	})

	router.HandleFunc("PUT /api/v1/storedata", api.StoredataHandler)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}
	done := make(chan os.Signal,1)
	slog.Info("Starting server on localhost:8080")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-done
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", slog.String("error", err.Error()))
	}

}
