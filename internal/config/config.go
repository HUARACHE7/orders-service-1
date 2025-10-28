package config

import "os"

const (
	DBUser   = "wb_user"
	DBPass   = "wb_pass"
	DBPort   = "5432"
	DBName   = "orders_db"
	
	NatsClusterID = "test-cluster" 
	NatsClientID  = "orders-service-client" 
	NatsChannel   = "orders-channel"
	NatsURL       = "nats://nats-streaming:4222"
	
	HTTPServerPort = ":8080"
)

func GetConnectionString() string {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = DBName
	}
	
	return "user=" + DBUser + 
		" password=" + DBPass + 
		" dbname=" + dbName +
		" host=" + dbHost + 
		" port=" + DBPort + 
		" sslmode=disable"
}