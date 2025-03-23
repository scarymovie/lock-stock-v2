package main

import (
	"context"
	"errors"
	"fmt"
	"lock-stock-v2/wire"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting server: %s\n", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

	wg.Wait()
	fmt.Println("Server exited.")
}
