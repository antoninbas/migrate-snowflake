package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/snowflake"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	sf "github.com/snowflakedb/gosnowflake"
	"go.uber.org/zap"
)

var (
	logger logr.Logger

	source    string
	database  string
	schema    string
	warehouse string
)

func runMigrations(ctx context.Context, db *sql.DB) error {
	dbInstance, err := snowflake.WithInstance(db, &snowflake.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(source, "snowflake", dbInstance)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database is already up-to-date")
			return nil
		}
		return fmt.Errorf("error when running migrations: %w", err)
	}
	return nil
}

func run(ctx context.Context) error {
	options := make([]func(*sf.Config), 0)
	if database != "" {
		options = append(options, SetDatabase(database))
	}
	if schema != "" {
		options = append(options, SetSchema(schema))
	}
	if warehouse != "" {
		options = append(options, SetWarehouse(warehouse))
	}
	dsn, _, err := GetDSN(options...)
	if err != nil {
		return fmt.Errorf("failed to create DSN: %w", err)
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to Snowflake: %w", err)
	}

	logger.Info("Migrating Snowflake database")
	if err := runMigrations(ctx, db); err != nil {
		return err
	}
	logger.Info("Migrated Snowflake database")

	return nil
}

func main() {
	flag.StringVar(&source, "source", "", "Source for migrations (only file://... supported)")
	flag.StringVar(&database, "database", "", "Database to migrate")
	flag.StringVar(&schema, "schema", "", "Schema to migrate")
	flag.StringVar(&warehouse, "warehouse", "", "Warehouse to use for queries")
	flag.Parse()

	if source == "" {
		log.Fatalf("-source is required")
	}

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zc.DisableStacktrace = true
	zapLog, err := zc.Build()
	if err != nil {
		panic("Cannot initialize Zap logger")
	}
	logger = zapr.NewLogger(zapLog)
	if err := run(context.Background()); err != nil {
		logger.Error(err, "Failed to migrate database")
		os.Exit(1)
	}
}
