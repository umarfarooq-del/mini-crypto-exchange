package models

import (
	"time"
)

// TradingPair represents a currency pair
type TradingPair struct {
	Base  string // e.g., "BTC"
	Quote string // e.g., "USDT"
}

// Order represents a user's buy/sell intent
type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Pair      string    `json:"pair"` // e.g., "BTC/USDT"
	Side      string    `json:"side"` // "buy" or "sell"
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Filled    float64   `json:"filled"`
	Status    string    `json:"status"` // "open", "partial", "filled"
	CreatedAt time.Time `json:"created_at"`
}

// Remaining returns the unfilled quantity
func (o *Order) Remaining() float64 {
	return o.Quantity - o.Filled
}

// Trade represents a matched trade
type Trade struct {
	ID          int64     `json:"id"`
	BuyOrderID  int64     `json:"buy_order_id"`
	SellOrderID int64     `json:"sell_order_id"`
	Pair        string    `json:"pair"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}
