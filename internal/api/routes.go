package api

import (
    "database/sql"
    "order-matching-system/internal/repository"
    "order-matching-system/internal/service"
    
    "github.com/gin-gonic/gin"
)

type Server struct {
    router       *gin.Engine
    handlers     *Handlers
    orderService *service.OrderService
}

func NewServer(matchingEngine *service.MatchingEngine, db *sql.DB) *Server {
    orderRepo := repository.NewOrderRepository(db)
    tradeRepo := repository.NewTradeRepository(db)
    orderService := service.NewOrderService(orderRepo, tradeRepo, matchingEngine)
    handlers := NewHandlers(orderService)
    
    router := gin.New()
    router.Use(LoggerMiddleware())
    router.Use(CORSMiddleware())
    router.Use(ErrorHandlerMiddleware())
    
    server := &Server{
        router:       router,
        handlers:     handlers,
        orderService: orderService,
    }
    
    server.setupRoutes()
    return server
}

func (s *Server) setupRoutes() {
    api := s.router.Group("/api/v1")
    
    // Health check
    api.GET("/health", s.handlers.Health)
    
    // Order operations
    api.POST("/orders", s.handlers.PlaceOrder)
    api.DELETE("/orders/:orderId", s.handlers.CancelOrder)
    api.GET("/orders/:orderId", s.handlers.GetOrder)
    
    // Market data
    api.GET("/orderbook", s.handlers.GetOrderBook)
    api.GET("/trades", s.handlers.GetTrades)
}

func (s *Server) Start(port string) error {
    return s.router.Run(":" + port)
}