package api

import (
    "net/http"
    "order-matching-system/internal/models"
    "order-matching-system/internal/service"
    "order-matching-system/internal/utils"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type Handlers struct {
    orderService *service.OrderService
}

func NewHandlers(orderService *service.OrderService) *Handlers {
    return &Handlers{
        orderService: orderService,
    }
}

func (h *Handlers) PlaceOrder(c *gin.Context) {
    var req models.PlaceOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, "Invalid request body")
        return
    }
    
    order, err := h.orderService.PlaceOrder(&req)
    if err != nil {
        utils.Error(c, err)
        return
    }
    
    utils.Success(c, order)
}

func (h *Handlers) CancelOrder(c *gin.Context) {
    orderID := c.Param("orderId")
    if orderID == "" {
        utils.BadRequest(c, "Order ID is required")
        return
    }
    
    err := h.orderService.CancelOrder(orderID)
    if err != nil {
        utils.Error(c, err)
        return
    }
    
    utils.Success(c, gin.H{"message": "Order canceled successfully"})
}

func (h *Handlers) GetOrder(c *gin.Context) {
    orderID := c.Param("orderId")
    if orderID == "" {
        utils.BadRequest(c, "Order ID is required")
        return
    }
    
    order, err := h.orderService.GetOrder(orderID)
    if err != nil {
        utils.Error(c, err)
        return
    }
    
    utils.Success(c, order)
}

func (h *Handlers) GetOrderBook(c *gin.Context) {
    symbol := c.Query("symbol")
    if symbol == "" {
        utils.BadRequest(c, "Symbol is required")
        return
    }
    
    orderBook := h.orderService.GetOrderBook(symbol)
    utils.Success(c, orderBook)
}

func (h *Handlers) GetTrades(c *gin.Context) {
    symbol := c.Query("symbol")
    if symbol == "" {
        utils.BadRequest(c, "Symbol is required")
        return
    }
    
    limit := 50
    if limitStr := c.Query("limit"); limitStr != "" {
        if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
            limit = parsedLimit
        }
    }
    
    trades, err := h.orderService.GetTrades(symbol, limit)
    if err != nil {
        utils.Error(c, err)
        return
    }
    
    utils.Success(c, trades)
}

func (h *Handlers) Health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "service": "order-matching-system",
    })
}