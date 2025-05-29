package models

import (
    "sort"
    
    "github.com/shopspring/decimal"
)

type OrderBook struct {
    Symbol string       `json:"symbol"`
    Bids   []PriceLevel `json:"bids"`
    Asks   []PriceLevel `json:"asks"`
}

type PriceLevel struct {
    Price    decimal.Decimal `json:"price"`
    Quantity decimal.Decimal `json:"quantity"`
    Orders   int             `json:"orders"`
}

func NewOrderBook(symbol string, orders []Order) *OrderBook {
    ob := &OrderBook{
        Symbol: symbol,
        Bids:   make([]PriceLevel, 0),
        Asks:   make([]PriceLevel, 0),
    }
    
    // Group orders by price level
    bidLevels := make(map[string]PriceLevel)
    askLevels := make(map[string]PriceLevel)
    
    for _, order := range orders {
        if order.Status != OPEN && order.Status != PARTIAL {
            continue
        }
        
        if order.Type == MARKET {
            continue // Market orders don't sit on the book
        }
        
        priceStr := order.Price.String()
        
        if order.Side == BUY {
            level := bidLevels[priceStr]
            level.Price = *order.Price
            level.Quantity = level.Quantity.Add(order.RemainingQuantity)
            level.Orders++
            bidLevels[priceStr] = level
        } else {
            level := askLevels[priceStr]
            level.Price = *order.Price
            level.Quantity = level.Quantity.Add(order.RemainingQuantity)
            level.Orders++
            askLevels[priceStr] = level
        }
    }
    
    // Convert maps to slices
    for _, level := range bidLevels {
        ob.Bids = append(ob.Bids, level)
    }
    for _, level := range askLevels {
        ob.Asks = append(ob.Asks, level)
    }
    
    // Sort bids (highest price first)
    sort.Slice(ob.Bids, func(i, j int) bool {
        return ob.Bids[i].Price.GreaterThan(ob.Bids[j].Price)
    })
    
    // Sort asks (lowest price first)
    sort.Slice(ob.Asks, func(i, j int) bool {
        return ob.Asks[i].Price.LessThan(ob.Asks[j].Price)
    })
    
    return ob
}

// Custom errors
var (
    ErrInvalidQuantity      = NewAPIError(400, "INVALID_QUANTITY", "Quantity must be positive")
    ErrInvalidPrice         = NewAPIError(400, "INVALID_PRICE", "Price must be positive for limit orders")
    ErrMarketOrderWithPrice = NewAPIError(400, "MARKET_ORDER_WITH_PRICE", "Market orders cannot have a price")
    ErrInvalidSide          = NewAPIError(400, "INVALID_SIDE", "Side must be 'buy' or 'sell'")
    ErrOrderNotFound        = NewAPIError(404, "ORDER_NOT_FOUND", "Order not found")
    ErrOrderAlreadyFilled   = NewAPIError(400, "ORDER_ALREADY_FILLED", "Order is already filled")
    ErrOrderAlreadyCanceled = NewAPIError(400, "ORDER_ALREADY_CANCELED", "Order is already canceled")
)

type APIError struct {
    Code    int    `json:"code"`
    Type    string `json:"type"`
    Message string `json:"message"`
}

func (e *APIError) Error() string {
    return e.Message
}

func NewAPIError(code int, errorType, message string) *APIError {
    return &APIError{
        Code:    code,
        Type:    errorType,
        Message: message,
    }
}