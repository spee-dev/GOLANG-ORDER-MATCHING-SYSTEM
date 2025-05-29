package service

import (
    "database/sql"
    "log"
    "order-matching-system/internal/models"
    "order-matching-system/internal/repository"
    "sync"
    "time"
    
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type MatchingEngine struct {
    db              *sql.DB
    orderRepo       *repository.OrderRepository
    tradeRepo       *repository.TradeRepository
    orderChannel    chan *models.Order
    cancelChannel   chan string
    orderBooks      map[string]*InMemoryOrderBook
    mutex           sync.RWMutex
}

type InMemoryOrderBook struct {
    Symbol string
    Bids   []*models.Order // Sorted by price (desc) then time (asc)
    Asks   []*models.Order // Sorted by price (asc) then time (asc)
    mutex  sync.RWMutex
}

func NewMatchingEngine(db *sql.DB) *MatchingEngine {
    return &MatchingEngine{
        db:            db,
        orderRepo:     repository.NewOrderRepository(db),
        tradeRepo:     repository.NewTradeRepository(db),
        orderChannel:  make(chan *models.Order, 1000),
        cancelChannel: make(chan string, 1000),
        orderBooks:    make(map[string]*InMemoryOrderBook),
    }
}

func (me *MatchingEngine) Start() {
    log.Println("Starting matching engine...")
    
    me.loadExistingOrders()
    
    for {
        select {
        case order := <-me.orderChannel:
            me.processOrder(order)
        case orderID := <-me.cancelChannel:
            me.processCancelOrder(orderID)
        }
    }
}

func (me *MatchingEngine) loadExistingOrders() {
   
    symbols := []string{"BTCUSD", "ETHUSD"}
    
    for _, symbol := range symbols {
        orders, err := me.orderRepo.GetOpenOrdersBySymbol(symbol)
        if err != nil {
            log.Printf("Error loading orders for %s: %v", symbol, err)
            continue
        }
        
        me.getOrCreateOrderBook(symbol)
        for _, order := range orders {
            me.addToOrderBook(&order)
        }
    }
}

func (me *MatchingEngine) PlaceOrder(order *models.Order) {
     me.processOrder(order)
    
}

func (me *MatchingEngine) CancelOrder(orderID string) {
    me.cancelChannel <- orderID
}

func (me *MatchingEngine) processOrder(order *models.Order) {
    log.Printf("Processing order: %s %s %s %v @ %v", order.ID, order.Side, order.Type, order.RemainingQuantity, order.Price)
    
    orderBook := me.getOrCreateOrderBook(order.Symbol)
    
    if order.Type == models.MARKET {
        me.processMarketOrder(order, orderBook)
    } else {
        me.processLimitOrder(order, orderBook)
    }
}

func (me *MatchingEngine) processMarketOrder(order *models.Order, orderBook *InMemoryOrderBook) {
    orderBook.mutex.Lock()
    defer orderBook.mutex.Unlock()
    
    tx, err := me.db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return
    }
    defer tx.Rollback()
    
    var oppositeOrders []*models.Order
    if order.Side == models.BUY {
        oppositeOrders = orderBook.Asks
    } else {
        oppositeOrders = orderBook.Bids
    }
    
    remainingQuantity := order.RemainingQuantity
    
    for i, restingOrder := range oppositeOrders {
        if remainingQuantity.IsZero() {
            break
        }
        
        matchQuantity := decimal.Min(remainingQuantity, restingOrder.RemainingQuantity)
        tradePrice := *restingOrder.Price 
        
        // Create trade
        trade := &models.Trade{
            ID:          uuid.New().String(),
            Symbol:      order.Symbol,
            Price:       tradePrice,
            Quantity:    matchQuantity,
            ExecutedAt:  time.Now(),
        }
        
        if order.Side == models.BUY {
            trade.BuyOrderID = order.ID
            trade.SellOrderID = restingOrder.ID
        } else {
            trade.BuyOrderID = restingOrder.ID
            trade.SellOrderID = order.ID
        }
        
        // Update quantities
        remainingQuantity = remainingQuantity.Sub(matchQuantity)
        restingOrder.RemainingQuantity = restingOrder.RemainingQuantity.Sub(matchQuantity)
        restingOrder.UpdatedAt = time.Now()
        
        // Update statuses
        if restingOrder.RemainingQuantity.IsZero() {
            restingOrder.Status = models.FILLED
        } else {
            restingOrder.Status = models.PARTIAL
        }
        
        // Save to database
        if err := me.tradeRepo.CreateWithTx(tx, trade); err != nil {
            log.Printf("Error saving trade: %v", err)
            return
        }
        
        if err := me.orderRepo.UpdateWithTx(tx, restingOrder); err != nil {
            log.Printf("Error updating resting order: %v", err)
            return
        }
        
        // Remove filled orders from order book
        if restingOrder.RemainingQuantity.IsZero() {
            if order.Side == models.BUY {
                orderBook.Asks = append(orderBook.Asks[:i], orderBook.Asks[i+1:]...)
            } else {
                orderBook.Bids = append(orderBook.Bids[:i], orderBook.Bids[i+1:]...)
            }
        }
        
        log.Printf("Trade executed: %v @ %v", matchQuantity, tradePrice)
    }
    
    // Update market order status
    if remainingQuantity.IsZero() {
        order.Status = models.FILLED
    } else {
        order.Status = models.CANCELED 
    }
    order.RemainingQuantity = decimal.Zero
    order.UpdatedAt = time.Now()
    
    if err := me.orderRepo.UpdateWithTx(tx, order); err != nil {
        log.Printf("Error updating market order: %v", err)
        return
    }
    
    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return
    }
    
    log.Printf("Market order processed: %s", order.ID)
}

func (me *MatchingEngine) processLimitOrder(order *models.Order, orderBook *InMemoryOrderBook) {
    orderBook.mutex.Lock()
    defer orderBook.mutex.Unlock()
    
    tx, err := me.db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return
    }
    defer tx.Rollback()
    
    var oppositeOrders []*models.Order
    if order.Side == models.BUY {
        oppositeOrders = orderBook.Asks
    } else {
        oppositeOrders = orderBook.Bids
    }
    
    remainingQuantity := order.RemainingQuantity
    matchedOrders := make([]*models.Order, 0)
    
    for _, restingOrder := range oppositeOrders {
        if remainingQuantity.IsZero() {
            break
        }
        
        // Check if prices cross
        var canMatch bool
        if order.Side == models.BUY {
            canMatch = order.Price.GreaterThanOrEqual(*restingOrder.Price)
        } else {
            canMatch = order.Price.LessThanOrEqual(*restingOrder.Price)
        }
        
        if !canMatch {
            break
        }
        
        matchQuantity := decimal.Min(remainingQuantity, restingOrder.RemainingQuantity)
        tradePrice := *restingOrder.Price 
        
        // Create trade
        trade := &models.Trade{
            ID:          uuid.New().String(),
            Symbol:      order.Symbol,
            Price:       tradePrice,
            Quantity:    matchQuantity,
            ExecutedAt:  time.Now(),
        }
        
        if order.Side == models.BUY {
            trade.BuyOrderID = order.ID
            trade.SellOrderID = restingOrder.ID
        } else {
            trade.BuyOrderID = restingOrder.ID
            trade.SellOrderID = order.ID
        }
        
        // Update quantities
        remainingQuantity = remainingQuantity.Sub(matchQuantity)
        restingOrder.RemainingQuantity = restingOrder.RemainingQuantity.Sub(matchQuantity)
        restingOrder.UpdatedAt = time.Now()
        
        // Update statuses
        if restingOrder.RemainingQuantity.IsZero() {
            restingOrder.Status = models.FILLED
        } else {
            restingOrder.Status = models.PARTIAL
        }
        
        // Save to database
        if err := me.tradeRepo.CreateWithTx(tx, trade); err != nil {
            log.Printf("Error saving trade: %v", err)
            return
        }
        
        if err := me.orderRepo.UpdateWithTx(tx, restingOrder); err != nil {
            log.Printf("Error updating resting order: %v", err)
            return
        }
        
        matchedOrders = append(matchedOrders, restingOrder)
        
        log.Printf("Trade executed: %v @ %v", matchQuantity, tradePrice)
    }
    
    // Remove filled orders from order book
    for _, matchedOrder := range matchedOrders {
        if matchedOrder.RemainingQuantity.IsZero() {
            me.removeFromOrderBook(orderBook, matchedOrder)
        }
    }
    
    // Update incoming order
    order.RemainingQuantity = remainingQuantity
    order.UpdatedAt = time.Now()
    
    if remainingQuantity.IsZero() {
        order.Status = models.FILLED
    } else if remainingQuantity.LessThan(order.InitialQuantity) {
        order.Status = models.PARTIAL
    } else {
        order.Status = models.OPEN
    }
    
    if err := me.orderRepo.UpdateWithTx(tx, order); err != nil {
        log.Printf("Error updating limit order: %v", err)
        return
    }
    
    // Add remaining quantity to order book
    if !remainingQuantity.IsZero() {
        me.addToOrderBook(order)
    }
    
    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return
    }
    
    log.Printf("Limit order processed: %s", order.ID)
}

func (me *MatchingEngine) processCancelOrder(orderID string) {
    log.Printf("Processing cancel order: %s", orderID)
    
    order, err := me.orderRepo.GetByID(orderID)
    if err != nil {
        log.Printf("Error getting order for cancellation: %v", err)
        return
    }
    
    if order.Status == models.FILLED || order.Status == models.CANCELED {
        log.Printf("Cannot cancel order %s: already %s", orderID, order.Status)
        return
    }
    
    orderBook := me.getOrCreateOrderBook(order.Symbol)
    orderBook.mutex.Lock()
    defer orderBook.mutex.Unlock()
    
    // Remove from order book
    me.removeFromOrderBook(orderBook, order)
    
    // Update order status
    order.Status = models.CANCELED
    order.UpdatedAt = time.Now()
    
    if err := me.orderRepo.Update(order); err != nil {
        log.Printf("Error updating canceled order: %v", err)
        return
    }
    
    log.Printf("Order canceled: %s", orderID)
}

func (me *MatchingEngine) getOrCreateOrderBook(symbol string) *InMemoryOrderBook {
    me.mutex.Lock()
    defer me.mutex.Unlock()
    
    if orderBook, exists := me.orderBooks[symbol]; exists {
        return orderBook
    }
    
    orderBook := &InMemoryOrderBook{
        Symbol: symbol,
        Bids:   make([]*models.Order, 0),
        Asks:   make([]*models.Order, 0),
    }
    
    me.orderBooks[symbol] = orderBook
    return orderBook
}

func (me *MatchingEngine) addToOrderBook(order *models.Order) {
    orderBook := me.getOrCreateOrderBook(order.Symbol)
    
    if order.Side == models.BUY {
        // Insert into bids (sorted by price desc, then time asc)
        insertIndex := 0
        for i, bid := range orderBook.Bids {
            if order.Price.GreaterThan(*bid.Price) ||
                (order.Price.Equal(*bid.Price) && order.CreatedAt.Before(bid.CreatedAt)) {
                insertIndex = i
                break
            }
            insertIndex = i + 1
        }
        orderBook.Bids = append(orderBook.Bids[:insertIndex], append([]*models.Order{order}, orderBook.Bids[insertIndex:]...)...)
    } else {
        // Insert into asks (sorted by price asc, then time asc)
        insertIndex := 0
        for i, ask := range orderBook.Asks {
            if order.Price.LessThan(*ask.Price) ||
                (order.Price.Equal(*ask.Price) && order.CreatedAt.Before(ask.CreatedAt)) {
                insertIndex = i
                break
            }
            insertIndex = i + 1
        }
        orderBook.Asks = append(orderBook.Asks[:insertIndex], append([]*models.Order{order}, orderBook.Asks[insertIndex:]...)...)
    }
}

func (me *MatchingEngine) removeFromOrderBook(orderBook *InMemoryOrderBook, order *models.Order) {
    if order.Side == models.BUY {
        for i, bid := range orderBook.Bids {
            if bid.ID == order.ID {
                orderBook.Bids = append(orderBook.Bids[:i], orderBook.Bids[i+1:]...)
                break
            }
        }
    } else {
        for i, ask := range orderBook.Asks {
            if ask.ID == order.ID {
                orderBook.Asks = append(orderBook.Asks[:i], orderBook.Asks[i+1:]...)
                break
            }
        }
    }
}

func (me *MatchingEngine) GetOrderBook(symbol string) *models.OrderBook {
    orders, err := me.orderRepo.GetOpenOrdersBySymbol(symbol)
    if err != nil {
        log.Printf("Error getting order book: %v", err)
        return &models.OrderBook{Symbol: symbol}
    }
    
    return models.NewOrderBook(symbol, orders)
}