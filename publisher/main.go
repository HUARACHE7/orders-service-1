package main

import (
	"io/ioutil"
	"log"

	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal("Ошибка подключения к NATS: ", err)
	}
	defer sc.Close()

	data, err := ioutil.ReadFile("../model.json")
	if err != nil {
		log.Fatalf("Не могу прочитать model.json: %v", err)
	}
	
	newUUID := uuid.New().String()
	msg := strings.Replace(string(data), "b563feb7b2b84b6test", newUUID, -1)
	msg = strings.Replace(msg, "WBILMTESTTRACK", newUUID[:12], -1) 

	err = sc.Publish("orders-channel", []byte(msg))
	if err != nil {
		log.Fatalf("Ошибка публикации: %v", err)
	}

	log.Printf("Опубликован заказ с UID: %s", newUUID)
}
