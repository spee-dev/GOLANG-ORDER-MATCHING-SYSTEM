-- Order Matching System Database Schema

CREATE DATABASE IF NOT EXISTS ordermatching;
USE ordermatching;

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    side ENUM('buy', 'sell') NOT NULL,
    type ENUM('limit', 'market') NOT NULL,
    price DECIMAL(15,8) NULL,
    initial_quantity DECIMAL(15,8) NOT NULL,
    remaining_quantity DECIMAL(15,8) NOT NULL,
    status ENUM('open', 'filled', 'canceled', 'partial') NOT NULL DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_symbol_side_status (symbol, side, status),
    INDEX idx_price_created (price, created_at),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- Trades table
CREATE TABLE IF NOT EXISTS trades (
    id VARCHAR(36) PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    buy_order_id VARCHAR(36) NOT NULL,
    sell_order_id VARCHAR(36) NOT NULL,
    price DECIMAL(15,8) NOT NULL,
    quantity DECIMAL(15,8) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (buy_order_id) REFERENCES orders(id),
    FOREIGN KEY (sell_order_id) REFERENCES orders(id),
    INDEX idx_symbol_executed (symbol, executed_at),
    INDEX idx_buy_order (buy_order_id),
    INDEX idx_sell_order (sell_order_id)
);