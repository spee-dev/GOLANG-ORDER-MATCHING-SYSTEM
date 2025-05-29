package models

import (
    "time"
    
    "github.com/shopspring/decimal"
)

type OrderSide string
type OrderType string
type OrderStatus string

const (
    BUY  OrderSide = "buy"
    SELL OrderSide = "sell"
)

const (
    LIMIT  OrderType = "limit"
    MARKET OrderType = "market"
)

const (
    OPEN     OrderStatus = "open"
    FILLED   OrderStatus = "filled"
    CANCELED OrderStatus = "canceled"
    PARTIAL  OrderStatus = "partial"
)

type Order struct {
    ID                string          `json:"id"`
    Symbol            string          `json:"symbol"`
    Side              OrderSide       `json:"side"`
    Type              OrderType       `json:"type"`
    Price             *decimal.Decimal `json:"price,omitempty"`
    InitialQuantity   decimal.Decimal `json:"initial_quantity"`
    RemainingQuantity decimal.Decimal `json:"remaining_quantity"`
    Status            OrderStatus     `json:"status"`
    CreatedAt         time.Time       `json:"created_at"`
    UpdatedAt         time.Time       `json:"updated_at"`
}

type PlaceOrderRequest struct {
    Symbol   string          `json:"symbol" binding:"required"`
    Side     OrderSide       `json:"side" binding:"required"`
    Type     OrderType       `json:"type" binding:"required"`
    Price    *decimal.Decimal `json:"price,omitempty"`
    Quantity decimal.Decimal `json:"quantity" binding:"required"`
}

func (r *PlaceOrderRequest) Validate() error {
    if r.Quantity.LessThanOrEqual(decimal.Zero) {
        return ErrInvalidQuantity
    }
    
    if r.Type == LIMIT && (r.Price == nil || r.Price.LessThanOrEqual(decimal.Zero)) {
        return ErrInvalidPrice
    }
    
    if r.Type == MARKET && r.Price != nil {
        return ErrMarketOrderWithPrice
    }
    
    if r.Side != BUY && r.Side != SELL {
        return ErrInvalidSide
    }
    
    return nil
}