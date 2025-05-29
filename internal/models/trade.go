package models

import (
    "time"
    
    "github.com/shopspring/decimal"
)

type Trade struct {
    ID          string          `json:"id"`
    Symbol      string          `json:"symbol"`
    BuyOrderID  string          `json:"buy_order_id"`
    SellOrderID string          `json:"sell_order_id"`
    Price       decimal.Decimal `json:"price"`
    Quantity    decimal.Decimal `json:"quantity"`
    ExecutedAt  time.Time       `json:"executed_at"`
}