package app

import (
	"log"

	"github.com/Nebezdar/tech_lvl0.git/internal/cache"
	"github.com/Nebezdar/tech_lvl0.git/internal/db"
	"github.com/Nebezdar/tech_lvl0.git/internal/messaging"
)

type OrderService struct {
	db    *db.Database
	cache *cache.Cache
	nats  *messaging.Nats
}

func NewOrderService(db *db.Database, cache *cache.Cache, nats *messaging.Nats) *OrderService {
	return &OrderService{
		db:    db,
		cache: cache,
		nats:  nats,
	}
}

func (s *OrderService) Start() {
	_, err := s.nats.Subscribe("orders", s.handleOrder)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *OrderService) Stop() {
	// Cleanup tasks, if needed
}

func (s *OrderService) handleOrder(msg *nats.Msg) {
	orderID := string(msg.Data)
	log.Println("Processing order:", orderID)

	// Save the data to PostgreSQL
	if err := s.db.SaveOrder(orderID); err != nil {
		log.Println("Error saving to DB:", err)
	}

	// Cache the data in memory
	s.cache.Set(orderID, []byte(orderID))
}

func (s *OrderService) GetOrderByID(orderID string) ([]byte, error) {
	// Check the cache first
	if data, ok := s.cache.Get(orderID); ok {
		return data, nil
	}

	// If not in cache, try to get from the database
	orderData, err := s.db.GetOrder(orderID)
	if err != nil {
		return nil, err
	}

	// Cache the data for future use
	s.cache.Set(orderID, []byte(orderData))

	return []byte(orderData), nil
}
