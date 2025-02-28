package data

import (
	"errors"
	"database/sql"
)

type AuthorRepository struct {
	dbTemplate *DBTemplate
}

func NewAuthorRepository(dbTemplate *DBTemplate) *AuthorRepository {
	return &AuthorRepository{
		dbTemplate: dbTemplate,
	}
}

func (repo *AuthorRepository) Create(author Author) (Author, error) {
	query := `
		INSERT INTO authors (first_name, last_name, bio)
		VALUES ($1, $2, $3) RETURNING id`
	id, err := ExecuteInsert(repo.dbTemplate, query, author.FirstName, author.LastName, author.Bio)
	if err != nil {
		return Author{}, err
	}
	author.ID = id
	return author, nil
}

func (repo *AuthorRepository) GetById(id int) (Author, error) {
	query := `
		SELECT id, first_name, last_name, bio
		FROM authors
		WHERE id = $1`
	author, err := QueryStruct[Author](repo.dbTemplate, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Author{}, errors.New("author not found")
		}
		return Author{}, err
	}
	return *author, nil
}

func (repo *AuthorRepository) Update(id int, updated Author) (Author, error) {
	query := `
		UPDATE authors SET first_name = $1, last_name = $2, bio = $3
		WHERE id = $4`
	_, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, updated.FirstName, updated.LastName, updated.Bio, id)
	if err != nil {
		return Author{}, err
	}
	updated.ID = id
	return updated, nil
}

func (repo *AuthorRepository) Delete(id int) error {
	query := `DELETE FROM authors WHERE id = $1`
	rowsAffected, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("author not found")
	}
	return nil
}

func (repo *AuthorRepository) GetAll() ([]Author, error) {
	query := `
		SELECT id, first_name, last_name, bio
		FROM authors`
	authors, err := QueryStructs[Author](repo.dbTemplate, query)
	if err != nil {
		return nil, err
	}
	return authors, nil
}
