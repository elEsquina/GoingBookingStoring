package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func GetOrdersInTimeRange(start, end time.Time, repo *OrderRepository) ([]Order, error) {
	query := `
		SELECT o.id, o.total_price, o.created_at, o.status,
		       c.id AS "customer.id", c.name AS "customer.name", c.email AS "customer.email"
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		WHERE o.created_at BETWEEN $1 AND $2`
	orders, err := QueryStructs[Order](repo.dbTemplate, query, start, end)
	if err != nil {
		return nil, err
	}

	// Fetch order items for each order
	for i := range orders {
		orderItemsQuery := `
			SELECT oi.book_id, b.title, oi.quantity
			FROM order_items oi
			JOIN books b ON oi.book_id = b.id
			WHERE oi.order_id = $1`
		items, err := QueryStructs[OrderItem](repo.dbTemplate, orderItemsQuery, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

func generateSalesReport(repo *OrderRepository) (SalesReport, error) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)

	orders, err := GetOrdersInTimeRange(start, now, repo)
	if err != nil {
		return SalesReport{}, err
	}

	report := SalesReport{
		Timestamp:    now,
		TotalRevenue: 0,
		TotalOrders:  len(orders),
	}

	bookSalesMap := make(map[int]*BookSales)
	for _, order := range orders {
		report.TotalRevenue += order.TotalPrice
		for _, item := range order.Items {
			if _, exists := bookSalesMap[item.Book.ID]; !exists {
				bookSalesMap[item.Book.ID] = &BookSales{
					Book:     item.Book,
					Quantity: 0,
				}
			}
			bookSalesMap[item.Book.ID].Quantity += item.Quantity
		}
	}

	for _, sales := range bookSalesMap {
		report.TopSellingBooks = append(report.TopSellingBooks, *sales)
	}

	return report, nil
}

func saveReport(report SalesReport) error {
	filename := fmt.Sprintf("output-reports/report_%s.json", report.Timestamp.Format("20060102150405"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(report)
}

func StartReportGenerator(store *DBTemplate) {
	dao, err := GetDAO[Order]("order", store)
	if err != nil {
		log.Fatalf("Failed to get DAO for order: %v", err)
	}
	repo, ok := dao.(*OrderRepository)
	if !ok {
		log.Fatalf("Failed to cast DAO to *OrderRepository")
	}

	ticker := time.NewTicker(24 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		report, err := generateSalesReport(repo)
		if err != nil {
			log.Println("Error generating sales report:", err)
			continue
		}
		if err := saveReport(report); err != nil {
			log.Println("Error saving sales report:", err)
		}
	}
}
