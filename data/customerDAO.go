package data

import (
	"errors"
	"database/sql"
)

type CustomerRepository struct {
	dbTemplate *DBTemplate
}

func NewCustomerRepository(dbTemplate *DBTemplate) *CustomerRepository {
	return &CustomerRepository{
		dbTemplate: dbTemplate,
	}
}

func (repo *CustomerRepository) Create(customer Customer) (Customer, error) {
	query := `
		INSERT INTO customers (name, email, address, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	id, err := ExecuteInsert(repo.dbTemplate, query, customer.Name, customer.Email, customer.Address, customer.CreatedAt)
	if err != nil {
		return Customer{}, err
	}
	customer.ID = id
	return customer, nil
}

func (repo *CustomerRepository) GetById(id int) (Customer, error) {
	query := `
		SELECT id, name, email, address, created_at
		FROM customers
		WHERE id = $1`
	customer, err := QueryStruct[Customer](repo.dbTemplate, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Customer{}, errors.New("customer not found")
		}
		return Customer{}, err
	}
	return *customer, nil
}

func (repo *CustomerRepository) Update(id int, updated Customer) (Customer, error) {
	query := `
		UPDATE customers SET name = $1, email = $2, address = $3, created_at = $4
		WHERE id = $5`
	_, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, updated.Name, updated.Email, updated.Address, updated.CreatedAt, id)
	if err != nil {
		return Customer{}, err
	}
	updated.ID = id
	return updated, nil
}

func (repo *CustomerRepository) Delete(id int) error {
	query := `DELETE FROM customers WHERE id = $1`
	rowsAffected, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("customer not found")
	}
	return nil
}

func (repo *CustomerRepository) GetAll() ([]Customer, error) {
	query := `
		SELECT id, name, email, address, created_at
		FROM customers`
	customers, err := QueryStructs[Customer](repo.dbTemplate, query)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
