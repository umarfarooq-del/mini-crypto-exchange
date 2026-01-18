package services

import (
	"context"
	"mini-crypto-exchange/internal/apperrors"
	"mini-crypto-exchange/internal/engine"
	"mini-crypto-exchange/internal/models"
	"mini-crypto-exchange/internal/util"
	"sync"
	"log"
)

// PlaceOrderService defines the interface for placing orders
type PlaceOrderService interface {
	ValidateRequest(ctx context.Context, userID int64, pair string, side string, price float64, quantity float64) []*util.Error
	ProcessRequest(ctx context.Context, userID int64, pair string, side string, price float64, quantity float64) (*models.Order, []*models.Trade, error)
}

var placeOrderSvcStruct PlaceOrderService
var placeOrderServiceOnce sync.Once

type placeOrderService struct {
	engine *engine.MatchingEngine
	config *util.RouterConfig
}

// InitPlaceOrderService initializes the place order service
func InitPlaceOrderService(matchingEngine *engine.MatchingEngine, config *util.RouterConfig) PlaceOrderService {
	placeOrderServiceOnce.Do(func() {
		placeOrderSvcStruct = &placeOrderService{engine: matchingEngine, config: config}
	})
	return placeOrderSvcStruct
}

// GetPlaceOrderService returns the singleton instance
func GetPlaceOrderService() PlaceOrderService {
	if placeOrderSvcStruct == nil {
		panic("PlaceOrderService not initialized")
	}
	return placeOrderSvcStruct
}

// ValidateRequest validates the order request
func (s *placeOrderService) ValidateRequest(ctx context.Context, userID int64, pair string, side string, price float64, quantity float64) []*util.Error {
	var validationErrors []*util.Error

	if userID <= 0 {
		
		log.Println("Invalid user ID")
		validationErrors = append(validationErrors, util.ServerToError(apperrors.ErrInvalidUserID))
	}

	if price <= 0 {
		log.Println("Invalid price")
		validationErrors = append(validationErrors, util.ServerToError(apperrors.ErrInvalidPrice))
	}

	if quantity <= 0 {
		log.Println("Invalid quantity")
		validationErrors = append(validationErrors, util.ServerToError(apperrors.ErrInvalidQuantity))
	}

	if side != "buy" && side != "sell" {
		log.Println("Invalid side")
		validationErrors = append(validationErrors, util.ServerToError(apperrors.ErrInvalidSide))
	}

	return validationErrors
}

// ProcessRequest processes the order placement
func (s *placeOrderService) ProcessRequest(ctx context.Context, userID int64, pair string, side string, price float64, quantity float64) (*models.Order, []*models.Trade, error) {
	ob := s.engine.GetOrderBook(pair)
	if ob == nil {
		return nil, nil, apperrors.ErrPairNotFound
	}

	order, trades, err := s.engine.PlaceOrder(userID, pair, side, price, quantity)
	return order, trades, err
}
