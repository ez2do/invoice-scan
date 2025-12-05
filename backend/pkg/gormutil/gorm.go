package gormutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"invoice-scan/backend/pkg/config"
	"invoice-scan/backend/pkg/log"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// parseLogLevel converts string log level to GORM logger.LogLevel
func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		log.Warnf("Unknown database log level '%s', defaulting to 'info'", level)
		return logger.Info
	}
}

// GormLogWriter implements the logger.Writer interface to integrate with our logging system
type GormLogWriter struct{}

// Printf implements the logger.Writer interface
func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	log.Infof("[GORM] "+format, args...)
}

// GormLogger is a custom logger that integrates with our logging system
type GormLogger struct {
	writer logger.Writer
	config logger.Config
}

// LogMode implements the logger.Interface
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level
	return &newLogger
}

// Info implements the logger.Interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Info {
		l.writer.Printf("[INFO] "+msg, data...)
	}
}

// Warn implements the logger.Interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Warn {
		l.writer.Printf("[WARN] "+msg, data...)
	}
}

// Error implements the logger.Interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Error {
		l.writer.Printf("[ERROR] "+msg, data...)
	}
}

// Trace implements the logger.Interface for SQL query logging
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.config.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.config.LogLevel >= logger.Error && (!l.config.IgnoreRecordNotFoundError || err != gorm.ErrRecordNotFound):
		l.writer.Printf("[ERROR] [%.3fms] [rows:%d] %s | ERROR: %v", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= logger.Warn:
		l.writer.Printf("[SLOW] [%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.config.LogLevel >= logger.Info:
		l.writer.Printf("[SQL] [%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

// configureSSL configures SSL for MySQL connection
func configureSSL() error {
	if !config.GetBool("database.ssl.enabled") {
		return nil
	}

	caFile := config.GetString("database.ssl.ca_file")
	if caFile == "" {
		return fmt.Errorf("SSL enabled but ca_file not specified")
	}

	// Read CA certificate from file
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate from file '%s': %v", caFile, err)
	}

	log.Infof("Using CA certificate from file: %s", caFile)

	// Create certificate pool and add CA
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: !config.GetBool("database.ssl.verify_server_cert"),
	}

	// Register TLS config with MySQL driver
	err = mysqlDriver.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to register TLS config: %v", err)
	}

	return nil
}

// buildDSNWithSSL modifies the DSN to include SSL parameters if needed
func buildDSNWithSSL(originalDSN string) (string, error) {
	if !config.GetBool("database.ssl.enabled") {
		return originalDSN, nil
	}

	// Parse the DSN
	parsedURL, err := url.Parse(originalDSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSN: %v", err)
	}

	// Get query parameters
	query := parsedURL.Query()

	// Add SSL parameters
	query.Set("tls", "custom")

	// Update the URL with new query parameters
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// OpenDBConnection initializes the database connection
func OpenDBConnection() *gorm.DB {
	once.Do(func() {
		dsn := config.GetString("database.dsn")
		if dsn == "" {
			panic("database.dsn is not set")
		}

		// Configure SSL if enabled
		if err := configureSSL(); err != nil {
			panic("failed to configure SSL: " + err.Error())
		}

		// Build DSN with SSL parameters if needed
		finalDSN, err := buildDSNWithSSL(dsn)
		if err != nil {
			panic("failed to build DSN with SSL: " + err.Error())
		}

		// Configure GORM logger
		gormConfig := &gorm.Config{}

		// Enable SQL logging if configured
		if config.GetBool("database.log_sql") {
			logLevel := parseLogLevel(config.GetString("database.log_level"))
			gormConfig.Logger = &GormLogger{
				writer: &GormLogWriter{},
				config: logger.Config{
					LogLevel:                  logLevel, // Use configurable log level
					IgnoreRecordNotFoundError: true,     // Don't log ErrRecordNotFound errors
					Colorful:                  true,
					SlowThreshold:             200 * time.Millisecond, // Log slow queries
				},
			}
		}

		db, err = gorm.Open(mysql.Open(finalDSN), gormConfig)
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
	})

	return db
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if db == nil {
		return OpenDBConnection()
	}
	return db
}
