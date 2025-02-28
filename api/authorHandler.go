package api

import (
	"encoding/json"
	"errors"
	"finalproject/data"
	"net/http"
	"strconv"
	"strings"
)

func getAuthorRepoFromFactory(w http.ResponseWriter, r *http.Request) (data.IDAO[data.Author], error) {
	store, ok := r.Context().Value("memoryStore").(*data.DBTemplate)
	if !ok || store == nil {
		http.Error(w, "Store not found in context", http.StatusInternalServerError)
		return nil, errors.New("store not found in context")
	}

	repo, err := data.GetDAO[data.Author]("author", store)
	if err != nil {
		http.Error(w, "Failed to retrieve author repository", http.StatusInternalServerError)
		return nil, err
	}
	return repo, nil
}

func GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	repo, err := getAuthorRepoFromFactory(w, r)
	if err != nil {
		return
	}

	authors, err := repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve authors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	repo, err := getAuthorRepoFromFactory(w, r)
	if err != nil {
		return
	}

	var author data.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdAuthor, err := repo.Create(author)
	if err != nil {
		http.Error(w, "Failed to create author", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdAuthor)
}

func GetAuthorById(w http.ResponseWriter, r *http.Request) {
	repo, err := getAuthorRepoFromFactory(w, r)
	if err != nil {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	author, err := repo.GetById(id)
	if err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func UpdateAuthorById(w http.ResponseWriter, r *http.Request) {
	repo, err := getAuthorRepoFromFactory(w, r)
	if err != nil {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	var updatedAuthor data.Author
	if err := json.NewDecoder(r.Body).Decode(&updatedAuthor); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	author, err := repo.Update(id, updatedAuthor)
	if err != nil {
		http.Error(w, "Failed to update author", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func DeleteAuthorById(w http.ResponseWriter, r *http.Request) {
	repo, err := getAuthorRepoFromFactory(w, r)
	if err != nil {
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/authors/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	if err := repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete author", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
