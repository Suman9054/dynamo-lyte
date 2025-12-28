package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type NodeMeta struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func main() {
	// ------------------ etcd client ------------------
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		slog.Error("Failed to connect to etcd", "error", err)
		return
	}
	defer cli.Close()

	slog.Info("Connected to etcd")


	node := NodeMeta{
		IP:   "127.0.0.1",
		Port: 3000,
	}

	value, _ := json.Marshal(node)

	ctxPut, cancelPut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPut()

	resp, err := cli.Put(ctxPut, "/nodes/node1", string(value))
	if err != nil {
		slog.Error("Failed to put key to etcd", "error", err)
		return
	}
	fmt.Println("etcd Put Revision:", resp.Header.Revision)

	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Metadata Service!"))
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: mux,
	}

	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Starting server on http://127.0.0.1:8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server forced shutdown", "error", err)
	}

	slog.Info("Server exited cleanly")
}
