package main

import (
    "log"
    "invoice-backend/config"
    "invoice-backend/database"
    "invoice-backend/handlers"
    "invoice-backend/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load config
    cfg := config.Load()

    // Connect to database
    if err := database.Connect(cfg.DatabaseURL); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize schema
    if err := database.InitSchema(); err != nil {
        log.Fatal("Failed to initialize schema:", err)
    }

    // Setup router
    r := gin.Default()

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Auth routes
    authHandler := &handlers.AuthHandler{JWTSecret: cfg.JWTSecret}
    r.POST("/auth/register", authHandler.Register)
    r.POST("/auth/login", authHandler.Login)

    // Protected routes
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
    {
        syncHandler := &handlers.SyncHandler{}
        protected.POST("/sync/push", syncHandler.Push)
        protected.POST("/sync/pull", syncHandler.Pull)
    }

    log.Printf("🚀 Server starting on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}