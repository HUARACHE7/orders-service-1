package streaming

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
	"github.com/HUARACHE7/orders-service-1/internal/cache"
	"github.com/HUARACHE7/orders-service-1/internal/config"
	"github.com/HUARACHE7/orders-service-1/internal/db"
	"github.com/HUARACHE7/orders-service-1/internal/model"
)

func StartSubscriber(database *sql.DB, appCache *cache.AppCache) (stan.Conn, error) {
	sc, err := stan.Connect(config.NatsClusterID, config.NatsClientID, stan.NatsURL(config.NatsURL))
	if err != nil {
		return nil, err
	}

	msgHandler := func(msg *stan.Msg) {
		log.Printf("Получено сообщение. Обработка...")
		
		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Ошибка Unmarshal JSON (пропускаем): %v. Сообщение: %s", err, string(msg.Data))
			msg.Ack() 
			return
		}

		if order.OrderUID == "" {
			log.Println("Ошибка валидации: order_uid пуст. Пропускаем.")
			msg.Ack() 
			return
		}
		
		if err := db.InsertOrder(database, order); err != nil {
			log.Printf("Ошибка записи заказа %s в БД. NACK: %v.", order.OrderUID, err)
			return
		}

		appCache.Set(order.OrderUID, order)
		log.Printf("Заказ %s успешно сохранен в БД и кэше.", order.OrderUID)

		msg.Ack()
	}

	_, err = sc.Subscribe(
		config.NatsChannel,
		msgHandler,
		stan.SetManualAckMode(), 
		stan.DurableName(config.NatsClientID),
		stan.AckWait(30e9),      
	)
	if err != nil {
		sc.Close()
		return nil, err
	}
	
	log.Println("Успешная подписка на NATS-Streaming канал:", config.NatsChannel)
	return sc, nil
}
