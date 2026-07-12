package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDB(dbSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
