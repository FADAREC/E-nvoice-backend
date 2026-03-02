package main

import (
    "log"
    "os" // Added this
    "invoice-backend/config"
    "invoice-backend/database"
    "invoice-backend/handlers"
    "invoice-backend/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load config
    cfg := config.Load()

    // 1. Set mode BEFORE creating the router
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Connect to database
    if err := database.Connect(cfg.DatabaseURL); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize schema
    if err := database.InitSchema(); err != nil {
        log.Fatal("Failed to initialize schema:", err)
    }

    // 2. Setup router
    r := gin.Default()

    // 3. Corrected Proxy Setting (Fixed syntax here)
    r.SetTrustedProxies(nil)

    // Optional: Add a root route to avoid 404 logs
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Invoice API is live"})
    })

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