package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/fastbyt3/diy-mssql-test/pkg/repository"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/microsoft/go-mssqldb/azuread"
	"go.uber.org/zap"
)

type Service struct {
	Db *sql.DB
}

var (
	database   = utils.GetEnvVarOrDefault("MSSQL_DB_NAME", "master")
	username   = utils.GetEnvVarOrDefault("MSSQL_DB_USERNAME", "sa")
	password   = utils.GetEnvVarOrDefault("MSSQL_DB_PASSWORD", "Passw0rd123")
	port       = utils.GetEnvVarOrDefault("MSSQL_DB_PORT", "1433")
	host       = utils.GetEnvVarOrDefault("MSSQL_DB_HOST", "localhost")
	dbInstance *Service
)

func New() Service {
	if dbInstance != nil {
		return *dbInstance
	}

	connStr := fmt.Sprintf("Data Source=%s,%s;Initial Catalog=%s;User ID=%s;Password=%s;Connect Timeout=30;Authentication=SqlPassword;Application Name=azdata;Command Timeout=30", host, port, database, username, password)
	db, err := sql.Open(azuread.DriverName, connStr)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("Connected to MSSQL DB", zap.String("host", host), zap.String("port", port), zap.String("db", database))

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbInstance = &Service{
		Db: db,
	}
	return *dbInstance
}

func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.Db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		utils.Logger.Fatal("Database is down.. T_T", zap.String("error", err.Error()))
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.Db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *Service) Close() error {
	utils.Logger.Info("Disconnected from database", zap.String("db-name", database))
	return s.Db.Close()
}

func (s *Service) UserRepo() repository.UserRepository {
	return repository.NewUserRepo(s.Db)
}
