package server

import (
	"mini-crypto-exchange/internal/engine"
	"mini-crypto-exchange/internal/services"
	"mini-crypto-exchange/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

// NewRouter ...
func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

// InitializeRouter ...
func (r *Router) InitializeRouter(routerConfig *util.RouterConfig) {
	r.initializeRoutes(routerConfig)
}

// initializeRoutes ...
func (r *Router) initializeRoutes(routerConfig *util.RouterConfig) {
	s := (*r).PathPrefix("").Subrouter()

	// Order matching routes
	s.HandleFunc("/api/pairs",
		CreatePairHandler(routerConfig.MatchingEngine.(*engine.MatchingEngine), routerConfig)).
		Methods(http.MethodOptions, http.MethodPost).
		Name("CreatePairAPI")

	s.HandleFunc("/api/orders",
		PlaceOrderHandler(services.GetPlaceOrderService(), routerConfig)).
		Methods(http.MethodOptions, http.MethodPost).
		Name("PlaceOrderAPI")

	s.HandleFunc("/api/orders",
		GetOrdersHandler(services.GetOrderBookService(), routerConfig)).
		Methods(http.MethodGet).
		Name("GetOrdersAPI")

	s.HandleFunc("/api/orderbook",
		OrderBookHandler(services.GetOrderBookService(), routerConfig)).
		Methods(http.MethodOptions, http.MethodGet).
		Name("OrderBookAPI")
}
