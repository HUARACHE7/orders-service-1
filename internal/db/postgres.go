package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/HUARACHE7/orders-service-1/internal/config" 
	"github.com/HUARACHE7/orders-service-1/internal/model"
)

func ConnectDB() (*sql.DB, error) {
	connStr := config.GetConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Успешное подключение к PostgreSQL")
	return db, nil
}

func InsertOrder(db *sql.DB, order model.Order) error {
	orderData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	_, err = db.Exec("INSERT INTO orders (order_uid, data) VALUES ($1, $2) ON CONFLICT (order_uid) DO NOTHING", order.OrderUID, orderData)
	if err != nil {
		return fmt.Errorf("ошибка вставки в БД: %w", err)
	}
	return nil
}

func LoadOrdersFromDB(db *sql.DB) (map[string]model.Order, error) {
	rows, err := db.Query("SELECT order_uid, data FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cache := make(map[string]model.Order)
	for rows.Next() {
		var order model.Order
		var orderData []byte
		var uid string

		if err := rows.Scan(&uid, &orderData); err != nil {
			log.Printf("Ошибка сканирования строки БД: %v", err)
			continue
		}
		
		if err := json.Unmarshal(orderData, &order); err != nil {
			log.Printf("Ошибка Unmarshal данных из БД: %v", err)
			continue
		}
		cache[uid] = order
	}
	log.Printf("Кэш восстановлен. Загружено %d заказов из БД.", len(cache))
	return cache, nil
}