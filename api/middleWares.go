package api

import (
	"context"
	"finalproject/data"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Key struct {
	UserID int
}

var tokenStore = make(map[Key]uuid.UUID)

func BooksRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetAllBooks(w, r)
	} else if r.Method == http.MethodPost {
		CreateBook(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}
func LoginUserRouter(db *data.DBTemplate, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Login(db, w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}
func SignUserRouter(db *data.DBTemplate, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		SignUp(db, w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}
func BooksPathParamRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetBookById(w, r)
	} else if r.Method == http.MethodPut {
		UpdateBookById(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteBookById(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func AuthorsRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetAllAuthors(w, r)
	} else if r.Method == http.MethodPost {
		CreateAuthor(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func AuthorsPathParamRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetAuthorById(w, r)
	} else if r.Method == http.MethodPut {
		UpdateAuthorById(w, r)
	} else if r.Method == http.MethodDelete {
		DeleteAuthorById(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func validateEmail(email string) bool {
	// Define the regex pattern for validating email format
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	if !strings.ContainsAny(password, "0123456789") {
		return false
	}

	if !strings.ContainsAny(password, "!@#$%^&*()_+{}[]|:;<>,.?/") {
		return false
	}

	return true
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		exists := false

		for _, t := range tokenStore {
			if t.String() == token {
				exists = true
				break
			}
		}

		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer file.Close()

		logger := log.New(file, "", log.LstdFlags)
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		statusCode := lrw.statusCode
		logger.Printf("%s - %s %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, http.StatusText(statusCode))
	})
}

func ContextGeneration(store *data.DBTemplate, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		ctx = context.WithValue(ctx, "memoryStore", store)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingResponseWriter is a wrapper around an http.ResponseWriter that keeps track of the status code written to it.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
