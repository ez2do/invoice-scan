package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"invoice-scan/backend/internal/adapters/repo"
	adapterstorage "invoice-scan/backend/internal/adapters/storage"
	"invoice-scan/backend/internal/handlers"
	"invoice-scan/backend/pkg/config"
	pkgextraction "invoice-scan/backend/pkg/extraction"

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

	uploadPath := config.GetStringWithDefaultValue("storage.upload_path", "./uploads")
	baseURL := config.GetStringWithDefaultValue("storage.base_url", "http://localhost:3001")

	fileStorage, err := adapterstorage.NewLocalStorage(uploadPath, baseURL)
	if err != nil {
		log.Fatalf("Failed to create file storage: %v", err)
	}

	router := gin.Default()

	corsOrigins := config.GetStringSlice("cors.origin")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = corsOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Cache-Control", "X-File-Name"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	corsConfig.AllowCredentials = false
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	router.Static("/uploads", uploadPath)
	router.StaticFile("/ssl/rootCA.pem", "./ssl/rootCA.pem")

	extractionService, err := pkgextraction.NewGeminiExtraction(geminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create extraction service: %v", err)
	}
	defer func() {
		if err := extractionService.Close(); err != nil {
			log.Printf("Error closing extraction service: %v", err)
		}
	}()

	extractHandler := handlers.NewExtractHandler(extractionService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceRepo, fileStorage, extractionService)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler)
		v1.POST("/extract", extractHandler.Extract)
		v1.POST("/invoices/upload", invoiceHandler.Upload)
		v1.GET("/invoices", invoiceHandler.List)
		v1.GET("/invoices/:id", invoiceHandler.GetByID)
		v1.DELETE("/invoices/:id", invoiceHandler.Delete)
	}

	host := config.GetStringWithDefaultValue("server.host", "localhost")
	port := config.GetStringWithDefaultValue("server.port", "3001")
	sslEnabled := config.GetBoolWithDefaultValue("server.ssl.enabled", false)
	sslCertFile := config.GetStringWithDefaultValue("server.ssl.cert_file", "./ssl/cert.pem")
	sslKeyFile := config.GetStringWithDefaultValue("server.ssl.key_file", "./ssl/key.pem")

	addr := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if sslEnabled {
			log.Printf("Server starting on https://%s", addr)
			if err := srv.ListenAndServeTLS(sslCertFile, sslKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			log.Printf("Server starting on http://%s", addr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Failed to start server: %v", err)
			}
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
