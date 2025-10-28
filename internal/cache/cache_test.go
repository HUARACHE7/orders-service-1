package cache

import (
	"testing"
	"time"
	"github.com/HUARACHE7/orders-service-1/internal/model"
)

var testOrder = model.Order{
	OrderUID: "test_uid_123",
	TrackNumber: "TESTTRACK",
	DateCreated: time.Now(),
}

func TestCacheSetAndGet(t *testing.T) {
	c := NewCache()

	c.Set(testOrder.OrderUID, testOrder)

	retrievedOrder, found := c.Get(testOrder.OrderUID)

	if !found {
		t.Errorf("Get() failed: order not found in cache")
	}
	
	if retrievedOrder.OrderUID != testOrder.OrderUID {
		t.Errorf("Get() returned wrong UID. Got: %s, Want: %s", retrievedOrder.OrderUID, testOrder.OrderUID)
	}

	_, notFound := c.Get("non_existent_uid")
	if notFound {
		t.Errorf("Get() failed: found non-existent key")
	}
}

func TestCacheLoad(t *testing.T) {
	c := NewCache()
	
	initialData := map[string]model.Order{
		"a": testOrder,
	}

	c.Load(initialData)

	if _, found := c.Get("a"); !found {
		t.Errorf("Load() failed: key 'a' not found after loading")
	}

	if len(c.items) != 1 {
		t.Errorf("Load() failed: incorrect number of items. Got %d, Want 1", len(c.items))
	}
}
