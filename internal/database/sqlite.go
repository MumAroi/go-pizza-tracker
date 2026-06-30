package database

import (
	"fmt"
	"pizza-tracker/internal/orders"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDB(dbSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.AutoMigrate(&orders.Order{}, &orders.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
