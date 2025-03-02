package data

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DBTemplate represents a database connection and template

// UserRepository represents the user repository that will interact with the database
type UserRepository struct {
	DbTemplate *DBTemplate
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(DbTemplate *DBTemplate) *UserRepository {
	return &UserRepository{
		DbTemplate: DbTemplate,
	}
}

// CreateUser checks if the user exists and then creates a new user
func (r *UserRepository) Create(user User) (User, error) {
	// Check if the user already exists
	exists := r.ExistUser(user.Email)
	if exists {
		return user, fmt.Errorf("user with email %s already exists", user.Email)
	}
	var err error
	// Insert the new user into the database
	id, err := ExecuteInsert(r.DbTemplate, "INSERT INTO users (useremail, password_hash) VALUES ($1, $2) RETURNING id", user.Email, user.PasswordHash)

	if err != nil {
		return user, fmt.Errorf("error inserting user: %v", err)
	}
	fmt.Println("Inserted user ID:", id)
	user.ID = id

	return user, nil
}

// ExistUser checks if a user already exists in the database
func (r *UserRepository) ExistUser(email string) bool {
	var exists bool
	err := r.DbTemplate.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE useremail=$1)", email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// CheckUserCredentials checks if the provided credentials are valid
func (r *UserRepository) CheckUserCredentials(email, password string) (bool, error) {

	var storedHash string
	if r.DbTemplate.db == nil || r.DbTemplate == nil {
		fmt.Println("Database connection not available")
		return false, nil
	}

	err := r.DbTemplate.db.QueryRow("SELECT password_hash FROM users WHERE useremail=$1", email).Scan(&storedHash)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("hrrllo")
		return false, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		fmt.Println("hllo")
		return false, fmt.Errorf("error fetching user data")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		fmt.Println("hrro")
		return false, fmt.Errorf("invalid credentials")
	}

	return true, nil
}

// GetById retrieves a user by their ID
func (r *UserRepository) GetById(id int) (User, error) {
	var user User
	err := r.DbTemplate.db.QueryRow("SELECT id, useremail, password_hash FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user with id %d not found", id)
		}
		return User{}, fmt.Errorf("error fetching user: %v", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (User, error) {
	var user User
	err := r.DbTemplate.db.QueryRow("SELECT id, useremail, password_hash FROM users WHERE useremail=$1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user with email %s not found", email)
		}
		return User{}, fmt.Errorf("error fetching user: %v", err)
	}
	return user, nil
}

// Update updates a user's email and password
func (r *UserRepository) Update(id int, email, password string) (User, error) {
	// Check if the user exists first
	user, err := r.GetById(id)
	if err != nil {
		return User{}, err
	}

	// Update the email if necessary
	if email != "" {
		user.Email = email
	}

	// Hash the new password
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("error hashing new password: %v", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Perform the update
	_, err = r.DbTemplate.db.Exec("UPDATE users SET useremail=$1, password_hash=$2 WHERE id=$3", user.Email, user.PasswordHash, id)
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %v", err)
	}

	return user, nil
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(id int) error {
	// Check if the user exists first
	_, err := r.GetById(id)
	if err != nil {
		return err
	}

	// Perform the deletion
	_, err = r.DbTemplate.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll() ([]User, error) {
	rows, err := r.DbTemplate.db.Query("SELECT id, useremail, password_hash FROM users")
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return users, nil
}

// User struct represents the structure of a user in the database
