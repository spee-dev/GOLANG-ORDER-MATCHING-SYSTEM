package utils

import (
    "net/http"
    "order-matching-system/internal/models"
    
    "github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Code    int    `json:"code"`
    Type    string `json:"type"`
    Message string `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
    })
}

func Error(c *gin.Context, err error) {
    if apiErr, ok := err.(*models.APIError); ok {
        c.JSON(apiErr.Code, Response{
            Success: false,
            Error: &ErrorInfo{
                Code:    apiErr.Code,
                Type:    apiErr.Type,
                Message: apiErr.Message,
            },
        })
        return
    }
    
    // Default to internal server error
    c.JSON(http.StatusInternalServerError, Response{
        Success: false,
        Error: &ErrorInfo{
            Code:    http.StatusInternalServerError,
            Type:    "INTERNAL_ERROR",
            Message: "Internal server error",
        },
    })
}

func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Success: false,
        Error: &ErrorInfo{
            Code:    http.StatusBadRequest,
            Type:    "BAD_REQUEST",
            Message: message,
        },
    })
}