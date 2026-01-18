package server

import (
	"encoding/json"
	"mini-crypto-exchange/internal/apperrors"
	"mini-crypto-exchange/internal/services"
	"mini-crypto-exchange/internal/util"
	"net/http"

	"log"
)

type PlaceOrderRequest struct {
	UserID   int64   `json:"user_id"`
	Pair     string  `json:"pair"`
	Side     string  `json:"side"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type PlaceOrderResponse struct {
	Order  interface{} `json:"order,omitempty"`
	Trades interface{} `json:"trades,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// PlaceOrderHandler handles POST /api/orders
func PlaceOrderHandler(service services.PlaceOrderService, config *util.RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		var req PlaceOrderRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode request: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlaceOrderResponse{Error: "Invalid request body"})
			return
		}

		// Validate request
		validationErrors := service.ValidateRequest(ctx, req.UserID, req.Pair, req.Side, req.Price, req.Quantity)
		if len(validationErrors) > 0 {
			log.Println("Validation failed")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlaceOrderResponse{Error: "Validation failed"})
			return
		}

		// Process request
		order, trades, err := service.ProcessRequest(ctx, req.UserID, req.Pair, req.Side, req.Price, req.Quantity)
		if err != nil {
			log.Printf("Failed to place order: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader( err.(*apperrors.ServerError).HTTPResponseCode)
			json.NewEncoder(w).Encode(PlaceOrderResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PlaceOrderResponse{
			Order:  order,
			Trades: trades,
		})
	}
}
