package server

import (
	"encoding/json"
	"log"
	"mini-crypto-exchange/internal/models"
	"mini-crypto-exchange/internal/services"
	"mini-crypto-exchange/internal/util"
	"net/http"
	"strconv"
)

type OrderBookResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type GetOrdersResponse struct {
	Orders []*models.Order `json:"orders,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// OrderBookHandler handles GET /api/orderbook
func OrderBookHandler(service services.OrderBookService, config *util.RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		pair := request.URL.Query().Get("pair")
		if pair == "" {
			log.Printf("Missing pair parameter")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(OrderBookResponse{Error: "Missing pair parameter"})
			return
		}

		depth := 10
		if depthStr := request.URL.Query().Get("depth"); depthStr != "" {
			if d, err := strconv.Atoi(depthStr); err == nil && d > 0 {
				depth = d
			}
		}

		data, err := service.GetOrderBook(ctx, pair, depth)
		if err != nil {
			log.Printf("Failed to get order book: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(OrderBookResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OrderBookResponse{Data: data})
	}
}

// GetOrdersHandler handles GET /api/orders?user_id=X
func GetOrdersHandler(service services.OrderBookService, config *util.RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		userIDStr := request.URL.Query().Get("user_id")
		if userIDStr == "" {
			log.Printf("Missing user_id parameter")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GetOrdersResponse{Error: "Missing user_id parameter"})
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			log.Printf("Invalid user_id parameter")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GetOrdersResponse{Error: "Invalid user_id parameter"})
			return
		}

		orders, err := service.GetOrdersByUser(ctx, userID)
		if err != nil {
			log.Printf("Failed to get orders: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(GetOrdersResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetOrdersResponse{Orders: orders})
	}
}
