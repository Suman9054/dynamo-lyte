package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	
)

func main() {
	router := http.NewServeMux()

	

	server:=http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}
  done:=make(chan os.Signal)
	slog.Info("Starting server on localhost:8080")
	go func(){
       if err:=server.ListenAndServe();err!=nil{
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