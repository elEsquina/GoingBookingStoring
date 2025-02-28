package data

import (
	"errors"
	"database/sql"
)

type OrderRepository struct {
	dbTemplate *DBTemplate
}

func NewOrderRepository(dbTemplate *DBTemplate) *OrderRepository {
	return &OrderRepository{
		dbTemplate: dbTemplate,
	}
}

func (repo *OrderRepository) Create(order Order) (Order, error) {
	query := `
		INSERT INTO orders (customer_id, total_price, created_at, status)
		VALUES ($1, $2, $3, $4) RETURNING id`
	id, err := ExecuteInsert(repo.dbTemplate, query, order.Customer.ID, order.TotalPrice, order.CreatedAt, order.Status)
	if err != nil {
		return Order{}, err
	}
	order.ID = id
	return order, nil
}

func (repo *OrderRepository) GetById(id int) (Order, error) {
	query := `
		SELECT o.id, o.total_price, o.created_at, o.status,
		       c.id AS "customer.id", c.name AS "customer.name", c.email AS "customer.email"
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		WHERE o.id = $1`
	order, err := QueryStruct[Order](repo.dbTemplate, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, errors.New("order not found")
		}
		return Order{}, err
	}
	return *order, nil
}

func (repo *OrderRepository) Update(id int, updated Order) (Order, error) {
	query := `
		UPDATE orders SET customer_id = $1, total_price = $2, created_at = $3, status = $4
		WHERE id = $5`
	_, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, updated.Customer.ID, updated.TotalPrice, updated.CreatedAt, updated.Status, id)
	if err != nil {
		return Order{}, err
	}
	updated.ID = id
	return updated, nil
}

func (repo *OrderRepository) Delete(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	rowsAffected, err := ExecuteUpdateOrDelete(repo.dbTemplate, query, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (repo *OrderRepository) GetByCustomerID(customerID int) ([]Order, error) {
	query := `
		SELECT o.id, o.total_price, o.created_at, o.status,
		       c.id AS "customer.id", c.name AS "customer.name", c.email AS "customer.email"
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		WHERE o.customer_id = $1`
	orders, err := QueryStructs[Order](repo.dbTemplate, query, customerID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (repo *OrderRepository) GetAll() ([]Order, error) {
	query := `
		SELECT o.id, o.total_price, o.created_at, o.status,
		       c.id AS "customer.id", c.name AS "customer.name", c.email AS "customer.email"
		FROM orders o
		JOIN customers c ON o.customer_id = c.id`
	orders, err := QueryStructs[Order](repo.dbTemplate, query)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
