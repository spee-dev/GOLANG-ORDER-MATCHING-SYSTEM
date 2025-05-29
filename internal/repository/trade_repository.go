package repository

import (
    "database/sql"
    "order-matching-system/internal/models"
)

type TradeRepository struct {
    db *sql.DB
}

func NewTradeRepository(db *sql.DB) *TradeRepository {
    return &TradeRepository{db: db}
}

func (r *TradeRepository) Create(trade *models.Trade) error {
    query := `
        INSERT INTO trades (id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := r.db.Exec(query,
        trade.ID,
        trade.Symbol,
        trade.BuyOrderID,
        trade.SellOrderID,
        trade.Price,
        trade.Quantity,
        trade.ExecutedAt,
    )
    
    return err
}

func (r *TradeRepository) CreateWithTx(tx *sql.Tx, trade *models.Trade) error {
    query := `
        INSERT INTO trades (id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := tx.Exec(query,
        trade.ID,
        trade.Symbol,
        trade.BuyOrderID,
        trade.SellOrderID,
        trade.Price,
        trade.Quantity,
        trade.ExecutedAt,
    )
    
    return err
}

func (r *TradeRepository) GetBySymbol(symbol string, limit int) ([]models.Trade, error) {
    query := `
        SELECT id, symbol, buy_order_id, sell_order_id, price, quantity, executed_at
        FROM trades
        WHERE symbol = ?
        ORDER BY executed_at DESC
        LIMIT ?
    `
    
    rows, err := r.db.Query(query, symbol, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var trades []models.Trade
    for rows.Next() {
        var trade models.Trade
        err := rows.Scan(
            &trade.ID,
            &trade.Symbol,
            &trade.BuyOrderID,
            &trade.SellOrderID,
            &trade.Price,
            &trade.Quantity,
            &trade.ExecutedAt,
        )
        if err != nil {
            return nil, err
        }
        trades = append(trades, trade)
    }
    
    return trades, nil
}