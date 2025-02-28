package data

import (
	"errors"
	"strings"
	"database/sql"
)

type BookRepository struct {
	dbTemplate *DBTemplate
}

func NewBookRepository(dbTemplate *DBTemplate) *BookRepository {
	return &BookRepository{
		dbTemplate: dbTemplate,
	}
}

func (repo *BookRepository) Create(book Book) (Book, error) {
	book.TextGenres = strings.Join(book.Genres, ",")
	query := `
		INSERT INTO books (title, author_id, genres, published_at, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	id, err := ExecuteInsert(repo.dbTemplate, query, book.Title, book.Author.ID, book.TextGenres, book.PublishedAt, book.Price, book.Stock)
	if err != nil {
		return Book{}, err
	}
	book.ID = id
	return book, nil
}

func (repo *BookRepository) GetById(id int) (Book, error) {
	query := `
		SELECT b.id, b.title, b.genres as "genres", b.published_at, b.price, b.stock,
		       a.id AS "author.id", a.first_name AS "author.first_name", a.last_name AS "author.last_name", a.bio AS "author.bio"
		FROM books b
		JOIN authors a ON b.author_id = a.id
		WHERE b.id = $1`
	book, err := QueryStruct[Book](repo.dbTemplate, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, errors.New("book not found")
		}
		return Book{}, err
	}
	book.Genres = strings.Split(book.TextGenres, ",")
	return *book, nil
}

func (repo *BookRepository) Update(id int, updated Book) (Book, error) {
	updated.TextGenres = strings.Join(updated.Genres, ",")
	query := `
		UPDATE books SET title = $1, author_id = $2, genres = $3, published_at = $4, price = $5, stock = $6
		WHERE id = $7`
	_, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, updated.Title, updated.Author.ID, updated.TextGenres, updated.PublishedAt, updated.Price, updated.Stock, id)
	if err != nil {
		return Book{}, err
	}
	updated.ID = id
	return updated, nil
}

func (repo *BookRepository) Delete(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	rowsAffected, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (repo *BookRepository) GetAll() ([]Book, error) {
	query := `
		SELECT b.id AS id, b.title, b.genres, b.published_at, b.price, b.stock,
		       a.id AS "author.id", a.first_name AS "author.first_name", a.last_name AS "author.last_name", a.bio AS "author.bio"
		FROM books b
		JOIN authors a ON b.author_id = a.id`
	books, err := QueryStructs[Book](repo.dbTemplate, query)
	if err != nil {
		return nil, err
	}
	for i := range books {
		books[i].Genres = strings.Split(books[i].TextGenres, ",")
	}
	return books, nil
}


func (repo *BookRepository) GetBookBySearchCriteria(s SearchCriteria) ([]Book, error) {
	query := `
		SELECT b.id AS id, b.title, b.genres, b.published_at, b.price, b.stock,
			a.id AS "author.id", a.first_name AS "author.first_name", a.last_name AS "author.last_name", a.bio AS "author.bio"
		FROM books b
		JOIN authors a ON b.author_id = a.id
		WHERE ($1 = '' OR b.title ILIKE $1)
		AND ($2 = '' OR a.first_name ILIKE $2)
		AND ($3 = '' OR b.genres ILIKE $3)
		ORDER BY b.title
	`

	books, err := QueryStructs[Book](repo.dbTemplate, query,
		s.Title,
		s.AuthorName, 
		"%" + s.Genre + "%") 


	if err != nil {
		return nil, err
	}
	for i := range books {
		books[i].Genres = strings.Split(books[i].TextGenres, ",")
	}
	return books, nil
}
