package services

import (
	"context"
	"mini-crypto-exchange/internal/engine"
	"mini-crypto-exchange/internal/models"
	"mini-crypto-exchange/internal/util"
	"sync"

	"log"
)

// OrderBookService defines the interface for querying order books
type OrderBookService interface {
	GetOrderBook(ctx context.Context, pair string, depth int) (map[string]interface{}, error)
	GetOrdersByUser(ctx context.Context, userID int64) ([]*models.Order, error)
}

var orderBookSvcStruct OrderBookService
var orderBookServiceOnce sync.Once

type orderBookService struct {
	engine *engine.MatchingEngine
	config *util.RouterConfig
}

// InitOrderBookService initializes the order book service
func InitOrderBookService(matchingEngine *engine.MatchingEngine, config *util.RouterConfig) OrderBookService {
	orderBookServiceOnce.Do(func() {
		orderBookSvcStruct = &orderBookService{engine: matchingEngine, config: config}
	})
	return orderBookSvcStruct
}

// GetOrderBookService returns the singleton instance
func GetOrderBookService() OrderBookService {
	if orderBookSvcStruct == nil {
		panic("OrderBookService not initialized")
	}
	return orderBookSvcStruct
}

// GetOrderBook returns the order book for a pair
func (s *orderBookService) GetOrderBook(ctx context.Context, pair string, depth int) (map[string]interface{}, error) {

	ob := s.engine.GetOrderBook(pair)
	if ob == nil {
		log.Printf("Order book not found for pair: %s", pair)
		return nil, nil
	}

	buys, sells := ob.GetDepth(depth)

	return map[string]interface{}{
		"pair": pair,
		"buy":  buys,
		"sell": sells,
	}, nil
}

// GetOrdersByUser returns all orders for a specific user
func (s *orderBookService) GetOrdersByUser(ctx context.Context, userID int64) ([]*models.Order, error) {

	orders := s.engine.GetOrdersByUser(userID)

	if orders == nil {
		orders = make([]*models.Order, 0)
	}

	return orders, nil
}
