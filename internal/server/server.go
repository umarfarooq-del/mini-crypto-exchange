package server

import (
	"log"
	"mini-crypto-exchange/internal/engine"
	"mini-crypto-exchange/internal/services"
	"mini-crypto-exchange/internal/util"
	"net"
	"net/http"
	"time"

	"github.com/soheilhy/cmux"
)

func RunServer() error {
	port := "50053"

	// Initialize matching engine
	matchingEngine := engine.NewMatchingEngine()
	routerConfigs := util.RouterConfig{
		MatchingEngine: matchingEngine,
	}

	// Initialize services
	services.InitPlaceOrderService(matchingEngine, &routerConfigs)
	services.InitOrderBookService(matchingEngine, &routerConfigs)

	// Setup router
	router := NewRouter()
	router.InitializeRouter(&routerConfigs)

	// Create listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// Create multiplexer
	mux := cmux.New(listener)

	// Setup servers
	httpServer := &http.Server{
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start servers
	go httpServer.Serve(mux.Match(cmux.HTTP1()))

	log.Printf("Starting HTTP and gRPC server on port %s", port)
	return mux.Serve()
}
