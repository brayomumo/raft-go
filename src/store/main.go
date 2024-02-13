package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hasKey := r.URL.Query().Has("key")
		if !hasKey{
			io.WriteString(w, "Home, where souls rest!")
			return
		}
		key := r.URL.Query().Get("key")
		if r.Method == "GET"{	
			success, str := store.HandleCommand("GET", []string{key,})
			if !success{
				io.WriteString(w, fmt.Sprintf("Error Getting key: %s", key))
				return
			}
			io.WriteString(w, str)
		}
		if r.Method == "POST"{
			hasValue := r.URL.Query().Has("value")
			if !hasValue{
				io.WriteString(w, "Malformed Request, No Value!")
				return
			}
			value := r.URL.Query().Get("value")
			success, str := store.HandleCommand("SET", []string{key, value})
			if !success{
				io.WriteString(w, "Error processing SET command!")
				return
			}
			io.WriteString(w, str)
		}
		
	})

	mux.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Key Endpoint Called!")
		io.WriteString(w, "Key endpoint called!")
	})
	// Add middlewares here 
	return mux
}


func main(){
	ctx := context.Background()
	defer ctx.Done()
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
		<- ctx.Done()
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
