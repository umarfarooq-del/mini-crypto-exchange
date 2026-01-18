package engine

import (
	"container/heap"
	"mini-crypto-exchange/internal/models"
	"sync"
)

// OrderBook manages buy and sell orders for a trading pair
type OrderBook struct {
	Pair        string
	BuyHeap     BuyHeap
	SellHeap    SellHeap
	mu          sync.Mutex
	nextTradeID int64
}

// NewOrderBook creates a new order book for a trading pair
func NewOrderBook(pair string) *OrderBook {
	return &OrderBook{
		Pair:        pair,
		BuyHeap:     make(BuyHeap, 0),
		SellHeap:    make(SellHeap, 0),
		nextTradeID: 1,
	}
}

// AddBuyOrder adds a buy order to the order book
func (ob *OrderBook) AddBuyOrder(order *models.Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	heap.Push(&ob.BuyHeap, order)
}

// AddSellOrder adds a sell order to the order book
func (ob *OrderBook) AddSellOrder(order *models.Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	heap.Push(&ob.SellHeap, order)
}

// GetBestBid returns the highest buy order without removing it
func (ob *OrderBook) GetBestBid() *models.Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if len(ob.BuyHeap) == 0 {
		return nil
	}
	return ob.BuyHeap[0]
}

// GetBestAsk returns the lowest sell order without removing it
func (ob *OrderBook) GetBestAsk() *models.Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if len(ob.SellHeap) == 0 {
		return nil
	}
	return ob.SellHeap[0]
}

// RemoveBestBid removes and returns the highest buy order
func (ob *OrderBook) RemoveBestBid() *models.Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if len(ob.BuyHeap) == 0 {
		return nil
	}
	return heap.Pop(&ob.BuyHeap).(*models.Order)
}

// RemoveBestAsk removes and returns the lowest sell order
func (ob *OrderBook) RemoveBestAsk() *models.Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if len(ob.SellHeap) == 0 {
		return nil
	}
	return heap.Pop(&ob.SellHeap).(*models.Order)
}

// GetDepth returns the order book depth (aggregated by price level)
func (ob *OrderBook) GetDepth(depth int) (buys []map[string]interface{}, sells []map[string]interface{}) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Aggregate buy orders by price
	buyMap := make(map[float64]float64)
	for _, order := range ob.BuyHeap {
		if order.Remaining() > 0 {
			buyMap[order.Price] += order.Remaining()
		}
	}

	// Aggregate sell orders by price
	sellMap := make(map[float64]float64)
	for _, order := range ob.SellHeap {
		if order.Remaining() > 0 {
			sellMap[order.Price] += order.Remaining()
		}
	}

	// Convert to sorted slices (simplified - in production use proper sorting)
	for price, qty := range buyMap {
		if len(buys) < depth {
			buys = append(buys, map[string]interface{}{
				"price":    price,
				"quantity": qty,
			})
		}
	}

	for price, qty := range sellMap {
		if len(sells) < depth {
			sells = append(sells, map[string]interface{}{
				"price":    price,
				"quantity": qty,
			})
		}
	}

	return buys, sells
}

// GetNextTradeID returns the next trade ID
func (ob *OrderBook) GetNextTradeID() int64 {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	id := ob.nextTradeID
	ob.nextTradeID++
	return id
}
