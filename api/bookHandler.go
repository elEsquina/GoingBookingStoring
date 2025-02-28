package api

import (
	"encoding/json"
	"errors"
	"finalproject/data"
	"net/http"
	"strconv"
	"strings"
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
	repo, err := getBookRepoFromFactory(w, r)
	if err != nil {
		return
	}

	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")

	var books []data.Book

	if title != "" || author != "" || genre != "" {
		searchCriteria := data.SearchCriteria{
			Title: title,
			AuthorName: author,
			Genre: genre,
		}
		books, err = repo.(*data.BookRepository).GetBookBySearchCriteria(searchCriteria)
	} else {
		books, err = repo.GetAll()
	}

	if err != nil {
		http.Error(w, "Failed to retrieve books" + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
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
		http.Error(w, "Failed to create book" + err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdBook)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
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

	book, err := repo.GetById(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
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
