package service

import (
    "order-matching-system/internal/models"
    "order-matching-system/internal/repository"
    "time"
    
    "github.com/google/uuid"
)

type OrderService struct {
    orderRepo      *repository.OrderRepository
    tradeRepo      *repository.TradeRepository
    matchingEngine *MatchingEngine
}

func NewOrderService(orderRepo *repository.OrderRepository, tradeRepo *repository.TradeRepository, matchingEngine *MatchingEngine) *OrderService {
    return &OrderService{
        orderRepo:      orderRepo,
        tradeRepo:      tradeRepo,
        matchingEngine: matchingEngine,
    }
}

func (s *OrderService) PlaceOrder(req *models.PlaceOrderRequest) (*models.Order, error) {
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    order := &models.Order{
        ID:                uuid.New().String(),
        Symbol:            req.Symbol,
        Side:              req.Side,
        Type:              req.Type,
        Price:             req.Price,
        InitialQuantity:   req.Quantity,
        RemainingQuantity: req.Quantity,
        Status:            models.OPEN,
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }
    
   
    if err := s.orderRepo.Create(order); err != nil {
        return nil, err
    }
    
    
   s.matchingEngine.PlaceOrder(order); 

     updatedOrder, err := s.orderRepo.GetByID(order.ID)
    if err != nil {
        return nil, err
    }
    return updatedOrder, nil
}

func (s *OrderService) CancelOrder(orderID string) error {
    order, err := s.orderRepo.GetByID(orderID)
    if err != nil {
        return err
    }
    
    if order.Status == models.FILLED {
        return models.ErrOrderAlreadyFilled
    }
    
    if order.Status == models.CANCELED {
        return models.ErrOrderAlreadyCanceled
    }
    
  
    s.matchingEngine.CancelOrder(orderID)
    
    return nil
}

func (s *OrderService) GetOrder(orderID string) (*models.Order, error) {
    return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) GetOrderBook(symbol string) *models.OrderBook {
    return s.matchingEngine.GetOrderBook(symbol)
}

func (s *OrderService) GetTrades(symbol string, limit int) ([]models.Trade, error) {
    if limit <= 0 || limit > 100 {
        limit = 50 
    }
    
    return s.tradeRepo.GetBySymbol(symbol, limit)
}