package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"invoice-scan/backend/internal/adapters/repo"
	"invoice-scan/backend/internal/handlers"
	"invoice-scan/backend/internal/services"
	"invoice-scan/backend/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	geminiAPIKey := config.GetStringWithDefaultValue("gemini.api_key", "")
	if geminiAPIKey == "" {
		log.Fatal("gemini.api_key environment variable is required")
	}

	dsn := getDSN()
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	invoiceRepo := repo.NewInvoiceGormRepo(gormDB)
	_ = invoiceRepo

	router := gin.Default()

	router.Use(gin.Recovery())

	corsOrigin := config.GetStringWithDefaultValue("cors.origin", "http://localhost:5173")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{corsOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	corsConfig.AllowCredentials = false
	router.Use(cors.New(corsConfig))

	extractionService, err := services.NewExtractionService(geminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create extraction service: %v", err)
	}
	defer extractionService.Close()

	extractHandler := handlers.NewExtractHandler(extractionService)

	api := router.Group("/api")
	{
		api.GET("/health", healthHandler)
		api.POST("/extract", extractHandler.Extract)
	}

	host := config.GetStringWithDefaultValue("server.host", "localhost")
	port := config.GetStringWithDefaultValue("server.port", "3001")

	addr := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getDSN() string {
	dbUser := config.GetStringWithDefaultValue("database.user", "root")
	dbPassword := config.GetStringWithDefaultValue("database.password", "")
	dbHost := config.GetStringWithDefaultValue("database.host", "localhost")
	dbPort := config.GetStringWithDefaultValue("database.port", "3306")
	dbName := config.GetStringWithDefaultValue("database.name", "invoice_scan")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
