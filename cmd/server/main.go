package main

import (
    "log"
    "order-matching-system/internal/api"
    "order-matching-system/internal/config"
    "order-matching-system/internal/database"
    "order-matching-system/internal/service"
)

func main() {
    
    cfg := config.Load()   // Load configuration
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // Run migrations
    if err := database.RunMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Initialize matching engine
    matchingEngine := service.NewMatchingEngine(db)
    
    // Start matching engine
    go matchingEngine.Start()
    server := api.NewServer(matchingEngine, db)
    log.Printf("Server starting on port %s", cfg.Port)
    if err := server.Start(cfg.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}