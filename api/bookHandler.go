package api

import (
	"encoding/json"
	"errors"
	"finalproject/caching"
	"finalproject/data"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getBookRepoFromFactory(w http.ResponseWriter, r *http.Request) (data.IDAO[data.Book], error) {
	store, ok := r.Context().Value("memoryStore").(*data.DBTemplate)
	if !ok || store == nil {
		http.Error(w, "Store not found in context", http.StatusInternalServerError)
		return nil, errors.New("store not found in context")
	}

	repo, err := data.GetDAO[data.Book]("book", store)
	if err != nil {
		http.Error(w, "Failed to retrieve book repository", http.StatusInternalServerError)
		return nil, err
	}
	return repo, nil
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")
	key := "books:" + title + genre + author

	if cachedData, exists := caching.Cache.Get(key); exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedData)
		return
	}

	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		http.Error(w, "Failed to get repository: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var books []data.Book
	if title != "" || author != "" || genre != "" {
		searchCriteria := data.SearchCriteria{Title: title, AuthorName: author, Genre: genre}
		books, err = repo.(*data.BookRepository).GetBookBySearchCriteria(searchCriteria)
	} else {
		books, err = repo.GetAll()
	}

	if err != nil {
		http.Error(w, "Failed to retrieve books: "+err.Error(), http.StatusInternalServerError)
		return
	}

	caching.Cache.Set(key, books, 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("book:%d", id)
	if cachedData, exists := caching.Cache.Get(key); exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedData)
		return
	}

	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		return
	}

	book, err := repo.GetById(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	caching.Cache.Set(key, book, 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		return
	}

	var book data.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdBook, err := repo.Create(book)
	if err != nil {
		http.Error(w, "Failed to create book"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdBook)
}

func UpdateBookById(w http.ResponseWriter, r *http.Request) {
	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook data.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	book, err := repo.Update(id, updatedBook)
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func DeleteBookById(w http.ResponseWriter, r *http.Request) {
	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
