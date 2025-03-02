package api

import (
	"encoding/json"
	"finalproject/data"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *data.DBTemplate, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	email, password := req.Email, req.Password
	fmt.Println(email, password)
	if !validateEmail(email) {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}
	if !validatePassword(password) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	repo := data.UserRepository{
		DbTemplate: db,
	}
	corr, erro := repo.CheckUserCredentials(email, password)

	if !corr {
		if erro.Error() == "invalid credentials" {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		} else if erro.Error() == fmt.Sprintf("user with email %s not found", email) {
			http.Error(w, fmt.Sprintf("user with email %s not found", email), http.StatusBadRequest)
			return
		} else {
			http.Error(w, "error fetching", http.StatusInternalServerError)
			return
		}
	}
	user, err := repo.GetByEmail(email)
	if err != nil {
		http.Error(w, "error fetching by email", http.StatusInternalServerError)
		return
	}
	key := Key{
		UserID: user.ID,
	}
	w.WriteHeader(http.StatusOK)
	newUUID, _ := uuid.NewUUID()
	tokenStore[key] = newUUID
	w.Write([]byte(newUUID.String()))

}

func SignUp(db *data.DBTemplate, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	email, password := req.Email, req.Password
	repo := data.UserRepository{
		DbTemplate: db,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating hash", http.StatusInternalServerError)
		return
	}
	user := data.User{
		ID:           0,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	user, err = repo.Create(user)
	if err != nil {
		http.Error(w, "Error creating User", http.StatusInternalServerError)
		return
	}
	newUUID, _ := uuid.NewUUID()
	tokenStore[Key{
		UserID: user.ID,
	}] = newUUID
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(newUUID.String()))

}
