package data

import (
	"time"
)

type Author struct {
	ID        int    `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Bio       string `json:"bio" db:"bio"`
}


type Book struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      Author    `json:"author" db:"author"`
	TextGenres  string    `json:"-" db:"genres"` 
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	Genres      []string  `json:"genres" db:"-"`
}

type Address struct {
	Street     string `json:"street" db:"street"`
	City       string `json:"city" db:"city"`
	State      string `json:"state" db:"state"`
	PostalCode string `json:"postal_code" db:"postal_code"`
	Country    string `json:"country" db:"country"`
}

type OrderItem struct {
	ID       int   `json:"id" db:"id"`
	Book     Book  `json:"book" db:"book"`
	Quantity int   `json:"quantity" db:"quantity"`
}

type Order struct {
	ID         int         `json:"id" db:"id"`
	Customer   Customer    `json:"customer" db:"customer"`
	Items      []OrderItem `json:"items" db:"items"`
	TotalPrice float64     `json:"total_price" db:"total_price"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	Status     string      `json:"status" db:"status"`
}

type Customer struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Address   Address   `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BookSales struct {
	Book     Book `json:"book" db:"book"`
	Quantity int  `json:"quantity_sold" db:"quantity_sold"`
}

type SalesReport struct {
	Timestamp       time.Time   `json:"timestamp" db:"timestamp"`
	TotalRevenue    float64     `json:"total_revenue" db:"total_revenue"`
	TotalOrders     int         `json:"total_orders" db:"total_orders"`
	TopSellingBooks []BookSales `json:"top_selling_books" db:"top_selling_books"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SearchCriteria struct {
	Title  string `json:"title"`
	AuthorName string `json:"author_name"`
	Genre  string `json:"genre"`
}