package main

import (
	"context"
	"finalproject/api"
	"finalproject/data"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/* 
TODO: 
	Postgres SQL: Done
	Greaceful shutdown: Done
	Logging: Done
	Authentication: Done
	Testing: Done
	Docker Container: Not Done
*/

func main() {
	template := data.NewDBTemplate("postgres://postgres:root@localhost:5432/finalproject?sslmode=disable")

	go data.StartReportGenerator(template)

	http.Handle("/login", api.RequestLogger( http.HandlerFunc(api.Login) ) )

	http.Handle("/books",
		api.RequestLogger(
			api.Authenticate(
				api.ContextGeneration(template, http.HandlerFunc(api.BooksRouter)),
			),
		),
	)

	http.Handle("/books/{id}",
		api.RequestLogger(
			api.Authenticate(
				api.ContextGeneration(template, http.HandlerFunc(api.BooksPathParamRouter)),
			),
		),
	)

	http.Handle("/authors",
		api.RequestLogger(
			api.Authenticate(
				api.ContextGeneration(template, http.HandlerFunc(api.AuthorsRouter)),
			),
		),
	)

	http.Handle("/authors/{id}",
		api.RequestLogger(
			api.Authenticate(
				api.ContextGeneration(template, http.HandlerFunc(api.AuthorsPathParamRouter)),
			),
		),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, 
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Server is starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
