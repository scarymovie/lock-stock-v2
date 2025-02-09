package main

import (
	"fmt"
	"lock-stock-v2/wire"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r, err := wire.InitializeRouter()
	if err != nil {
		fmt.Printf("Failed to initialize router: %v\n", err)
		return
	}

	port := ":8080"
	fmt.Printf("Server is running on %s\n", port)
	server := &http.Server{Addr: port, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	server.Close()
}
