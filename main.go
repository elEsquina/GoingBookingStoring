package main

import (
	"context"
	"finalproject/api"
	"finalproject/caching"
	"finalproject/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	template := data.NewDBTemplate("postgres://postgres:12345@localhost:5432/finalproject?sslmode=disable")

	go data.StartReportGenerator(template)

	// Initialize in-memory rate limiter using go-cache
	rateLimiter := caching.NewRateLimiter(1, 10, 1000*time.Second) // Adjust parameters as needed
	fmt.Println(rateLimiter)

	// Define routes with rate limiting
	http.Handle("/login", rateLimiter.Middleware(
		api.RequestLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api.LoginUserRouter(template, w, r)
		})),
	))

	http.Handle("/signup", rateLimiter.Middleware(
		api.RequestLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api.SignUserRouter(template, w, r)
		})),
	))

	http.Handle("/books",
		rateLimiter.Middleware(
			api.RequestLogger(
				api.Authenticate(
					api.ContextGeneration(template, http.HandlerFunc(api.BooksRouter)),
				),
			),
		),
	)

	http.Handle("/books/{id}",
		rateLimiter.Middleware(
			api.RequestLogger(
				api.Authenticate(
					api.ContextGeneration(template, http.HandlerFunc(api.BooksPathParamRouter)),
				),
			),
		),
	)

	http.Handle("/authors",
		rateLimiter.Middleware(
			api.RequestLogger(
				api.Authenticate(
					api.ContextGeneration(template, http.HandlerFunc(api.AuthorsRouter)),
				),
			),
		),
	)

	http.Handle("/authors/{id}",
		rateLimiter.Middleware(
			api.RequestLogger(
				api.Authenticate(
					api.ContextGeneration(template, http.HandlerFunc(api.AuthorsPathParamRouter)),
				),
			),
		),
	)

	// Start the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	// Graceful Shutdown Handling
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
