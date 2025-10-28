package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HUARACHE7/orders-service-1/internal/cache"
	"github.com/HUARACHE7/orders-service-1/internal/db"
	"github.com/HUARACHE7/orders-service-1/internal/http"
	"github.com/HUARACHE7/orders-service-1/internal/streaming"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	database, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer database.Close()

	appCache := cache.NewCache()
	restoredOrders, err := db.LoadOrdersFromDB(database)
	if err != nil {
		log.Fatalf("Не удалось восстановить кэш из БД: %v", err)
	}
	appCache.Load(restoredOrders)

	stanConn, err := streaming.StartSubscriber(database, appCache)
	if err != nil {
		log.Fatalf("Не удалось подключиться к NATS: %v", err)
	}
	defer stanConn.Close()

	go http.StartServer(appCache)

	log.Println("Сервис запущен. Нажмите Ctrl+C для выхода.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Сервис останавливается...")
}
