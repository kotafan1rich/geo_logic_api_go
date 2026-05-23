package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DB interface {
	GORM() *gorm.DB
}

type postgresDB struct {
	gormDB *gorm.DB
}

func (p *postgresDB) GORM() *gorm.DB {
	return p.gormDB
}

func New(dsn string, maxIdleConns, maxOpenConns int) (DB, error) {
	slog.Info("connecting to db")
	if dsn == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	gormLogger := gormLogger.Default.LogMode(gormLogger.Info)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db := &postgresDB{gormDB: dbConn}

	sqlDB, err := db.GORM().DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := dbConn.Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error; err != nil {
		return nil, err
	}

	slog.Info("connected succesfully")

	return db, nil
}
