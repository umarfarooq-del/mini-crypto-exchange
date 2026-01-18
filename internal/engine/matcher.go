package engine

import (
	"errors"
	"mini-crypto-exchange/internal/apperrors"
	"mini-crypto-exchange/internal/models"
	"sync"
	"time"
)

var (
	// ErrPairNotFound is returned when a trading pair is not found
	ErrPairNotFound = errors.New("trading pair not found")
)

// MatchingEngine manages order books and matching logic
type MatchingEngine struct {
	orderBooks  map[string]*OrderBook
	mu          sync.RWMutex
	nextOrderID int64
	orders      map[int64]*models.Order
	trades      []*models.Trade
}

// NewMatchingEngine creates a new matching engine
func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		orderBooks:  make(map[string]*OrderBook),
		nextOrderID: 1,
		orders:      make(map[int64]*models.Order),
		trades:      make([]*models.Trade, 0),
	}
}

// CreatePair creates a new trading pair
func (me *MatchingEngine) CreatePair(pair string) {
	me.mu.Lock()
	defer me.mu.Unlock()
	if _, exists := me.orderBooks[pair]; !exists {
		me.orderBooks[pair] = NewOrderBook(pair)
	}
}

// GetOrderBook returns the order book for a pair
func (me *MatchingEngine) GetOrderBook(pair string) *OrderBook {
	me.mu.RLock()
	defer me.mu.RUnlock()
	return me.orderBooks[pair]
}

// PlaceOrder places an order and attempts to match it
func (me *MatchingEngine) PlaceOrder(userID int64, pair string, side string, price float64, quantity float64) (*models.Order, []*models.Trade, error) {
	me.mu.Lock()
	ob, exists := me.orderBooks[pair]
	me.mu.Unlock()

	if !exists {
		return nil, nil, apperrors.ErrPairNotFound
	}

	// Create order
	order := &models.Order{
		ID:        me.getNextOrderID(),
		UserID:    userID,
		Pair:      pair,
		Side:      side,
		Price:     price,
		Quantity:  quantity,
		Filled:    0,
		Status:    "open",
		CreatedAt: time.Now(),
	}

	// Match order
	trades := me.matchOrder(ob, order)

	me.mu.Lock()
	me.orders[order.ID] = order
	me.mu.Unlock()

	// Update order status
	if order.Filled == order.Quantity {
		order.Status = "filled"
	} else if order.Filled > 0 {
		order.Status = "partial"
	}

	// Add remaining order to book if not fully filled
	if order.Remaining() > 0 {
		if side == "buy" {
			ob.AddBuyOrder(order)
		} else {
			ob.AddSellOrder(order)
		}
	}

	return order, trades, nil
}

// matchOrder matches an incoming order against the order book
func (me *MatchingEngine) matchOrder(ob *OrderBook, incomingOrder *models.Order) []*models.Trade {
	trades := make([]*models.Trade, 0)

	if incomingOrder.Side == "buy" {
		// Match buy order against sell orders
		for incomingOrder.Remaining() > 0 {
			bestAsk := ob.GetBestAsk()
			if bestAsk == nil || bestAsk.Price > incomingOrder.Price {
				break
			}

			// Match quantity
			matchQty := incomingOrder.Remaining()
			if bestAsk.Remaining() < matchQty {
				matchQty = bestAsk.Remaining()
			}

			// Update filled amounts
			incomingOrder.Filled += matchQty
			bestAsk.Filled += matchQty

			// Create trade
			trade := &models.Trade{
				ID:          ob.GetNextTradeID(),
				BuyOrderID:  incomingOrder.ID,
				SellOrderID: bestAsk.ID,
				Pair:        ob.Pair,
				Price:       bestAsk.Price,
				Quantity:    matchQty,
				CreatedAt:   time.Now(),
			}
			trades = append(trades, trade)
			me.trades = append(me.trades, trade)

			// Remove filled sell order
			if bestAsk.Remaining() == 0 {
				ob.RemoveBestAsk()
				bestAsk.Status = "filled"
			} else {
				bestAsk.Status = "partial"
			}
		}
	} else {
		// Match sell order against buy orders
		for incomingOrder.Remaining() > 0 {
			bestBid := ob.GetBestBid()
			if bestBid == nil || bestBid.Price < incomingOrder.Price {
				break
			}

			// Match quantity
			matchQty := incomingOrder.Remaining()
			if bestBid.Remaining() < matchQty {
				matchQty = bestBid.Remaining()
			}

			// Update filled amounts
			incomingOrder.Filled += matchQty
			bestBid.Filled += matchQty

			// Create trade
			trade := &models.Trade{
				ID:          ob.GetNextTradeID(),
				BuyOrderID:  bestBid.ID,
				SellOrderID: incomingOrder.ID,
				Pair:        ob.Pair,
				Price:       bestBid.Price,
				Quantity:    matchQty,
				CreatedAt:   time.Now(),
			}
			trades = append(trades, trade)
			me.trades = append(me.trades, trade)

			// Remove filled buy order
			if bestBid.Remaining() == 0 {
				ob.RemoveBestBid()
				bestBid.Status = "filled"
			} else {
				bestBid.Status = "partial"
			}
		}
	}

	return trades
}

func (me *MatchingEngine) getNextOrderID() int64 {
	me.mu.Lock()
	defer me.mu.Unlock()
	id := me.nextOrderID
	me.nextOrderID++
	return id
}

func (me *MatchingEngine) GetTrades() []*models.Trade {
	me.mu.RLock()
	defer me.mu.RUnlock()
	return me.trades
}

// GetOrdersByUser returns all orders for a specific user across all trading pairs
func (me *MatchingEngine) GetOrdersByUser(userID int64) []*models.Order {
	me.mu.RLock()
	defer me.mu.RUnlock()

	var userOrders []*models.Order

	for _, o := range me.orders {
		if o.UserID == userID {
			userOrders = append(userOrders, o)
		}
	}

	return userOrders
}
