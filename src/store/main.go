package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)


var (
	Address = "localhost"
	Port = "8000"
)


func NewHTTPServer(logger *logrus.Logger, store *Store) http.Handler{
	mux := http.NewServeMux()
	// TODO: Add Endpoints for set, get, delete, getAll, DeleteAll
	// addRoutes(
	// 	mux,
	// 	logger,
	// 	store,
	// )
	// // Add middlewares here 
	return mux
}


func main(){
	fmt.Println("KV Store up")

	store := NewStore()
	wg := &sync.WaitGroup{}

	srv := NewHTTPServer(
		logrus.New(),
		store,
	)

	httpServer := &http.Server{
		Addr: net.JoinHostPort(Address, Port),
		Handler: srv,
	}

	// kaskizie uko mbele
	go func(){
		logrus.Printf("Listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			logrus.Fatalf("error listening and serving: %s\n", err)
		}

	}()
	wg.Add(1)

	// Graceful shutdown
	go func(){
		defer wg.Done()
		// <- ctx.Done()
		// New context for the shutdown
		offContext := context.Background()
		offContext, cancel := context.WithTimeout(ctx, 10 *time.Second)
		defer cancel()

		if err := httpServer.Shutdown(offContext); err != nil {
			logrus.Fatal("error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

}
