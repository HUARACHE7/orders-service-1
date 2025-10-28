package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/HUARACHE7/orders-service-1/internal/model"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	
	testDB, err = ConnectDB() 
	if err != nil {
		log.Fatalf("Не удалось подключиться к тестовой БД (Postgres должен быть запущен в Docker): %v", err)
	}
	defer testDB.Close()
	
	testDB.Exec("DELETE FROM orders") 

	code := m.Run()
	
	testDB.Exec("DELETE FROM orders")
	
	os.Exit(code)
}

func createTestOrder(uid string) model.Order {
	return model.Order{
		OrderUID: uid,
		TrackNumber: "TEST",
		DateCreated: time.Now(),
		Delivery: model.Delivery{Name: "Test Name"},
		Payment: model.Payment{Transaction: "Test Txn"},
	}
}

func TestInsertAndLoad(t *testing.T) {
	testDB.Exec("DELETE FROM orders") 
	
	uid1 := "test_db_1"
	order1 := createTestOrder(uid1)
	
	if err := InsertOrder(testDB, order1); err != nil {
		t.Fatalf("InsertOrder failed: %v", err)
	}

	loadedCache, err := LoadOrdersFromDB(testDB)
	if err != nil {
		t.Fatalf("LoadOrdersFromDB failed: %v", err)
	}

	if len(loadedCache) != 1 {
		t.Fatalf("LoadOrdersFromDB returned incorrect count. Got %d, Want 1", len(loadedCache))
	}
	
	retrievedOrder, found := loadedCache[uid1]
	if !found || retrievedOrder.OrderUID != uid1 {
		t.Errorf("Loaded order mismatch. Found: %t, UID: %s", found, retrievedOrder.OrderUID)
	}

	order1.TrackNumber = "CHANGED_TRACK" 
	if err := InsertOrder(testDB, order1); err != nil {
		t.Fatalf("Second InsertOrder failed: %v", err)
	}
	
	loadedCacheAfterConflict, _ := LoadOrdersFromDB(testDB)
	if loadedCacheAfterConflict[uid1].TrackNumber != "TEST" {
		t.Errorf("ON CONFLICT failed: data was updated, expected no change.")
	}
}