package repository

import (
    "database/sql"
    "order-matching-system/internal/models"
    "log"
    "github.com/shopspring/decimal"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
    query := `
        INSERT INTO orders (id, symbol, side, type, price, initial_quantity, remaining_quantity, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    var price interface{} = nil
    if order.Price != nil {
        price = order.Price.String() // string representation of decimal
    }

    // Convert decimal.Decimal to string for quantities (assuming InitialQuantity and RemainingQuantity are decimal.Decimal)
    initialQty := order.InitialQuantity.String()
    remainingQty := order.RemainingQuantity.String()

    // Execute query passing actual values (not pointers)
    _, err := r.db.Exec(query,
        order.ID,
        order.Symbol,
        order.Side,
        order.Type,
        price,
        initialQty,
        remainingQty,
        order.Status,
        order.CreatedAt,
        order.UpdatedAt,
    )

    if err != nil {
        log.Printf("Failed to insert order: %v", err)
    }

    return err
}


    


func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
    query := `
        SELECT id, symbol, side, type, price, initial_quantity, remaining_quantity, status, created_at, updated_at
        FROM orders
        WHERE id = ?
    `
    
    row := r.db.QueryRow(query, id)
    return r.scanOrder(row)
}

func (r *OrderRepository) GetOpenOrdersBySymbol(symbol string) ([]models.Order, error) {
    query := `
        SELECT id, symbol, side, type, price, initial_quantity, remaining_quantity, status, created_at, updated_at
        FROM orders
        WHERE symbol = ? AND status IN ('open', 'partial')
        ORDER BY created_at ASC
    `
    
    rows, err := r.db.Query(query, symbol)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var orders []models.Order
    for rows.Next() {
        order, err := r.scanOrder(rows)
        if err != nil {
            return nil, err
        }
        orders = append(orders, *order)
    }
    
    return orders, nil
}

func (r *OrderRepository) Update(order *models.Order) error {
    query := `
        UPDATE orders
        SET remaining_quantity = ?, status = ?, updated_at = ?
        WHERE id = ?
    `
    
    _, err := r.db.Exec(query,
        order.RemainingQuantity,
        order.Status,
        order.UpdatedAt,
        order.ID,
    )
    
    return err
}

func (r *OrderRepository) UpdateWithTx(tx *sql.Tx, order *models.Order) error {
    query := `
        UPDATE orders
        SET remaining_quantity = ?, status = ?, updated_at = ?
        WHERE id = ?
    `
    
    _, err := tx.Exec(query,
        order.RemainingQuantity,
        order.Status,
        order.UpdatedAt,
        order.ID,
    )
    
    return err
}

func (r *OrderRepository) scanOrder(scanner interface {
    Scan(dest ...interface{}) error
}) (*models.Order, error) {
    var order models.Order
    var price sql.NullString
    
    err := scanner.Scan(
        &order.ID,
        &order.Symbol,
        &order.Side,
        &order.Type,
        &price,
        &order.InitialQuantity,
        &order.RemainingQuantity,
        &order.Status,
        &order.CreatedAt,
        &order.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, models.ErrOrderNotFound
        }
        return nil, err
    }
    
    if price.Valid {
        p, err := decimal.NewFromString(price.String)
        if err != nil {
            return nil, err
        }
        order.Price = &p
    }
    
    return &order, nil
}