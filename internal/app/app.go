package app

import (
	"fmt"
	"log/slog"
	"pizza-tracker/internal/database"
	"pizza-tracker/internal/order"

	"gorm.io/gorm"
)

type App struct {
	DB        *gorm.DB
	OrderRepo order.OrderRepository
}

func NewApp(dbPath string) (*App, error) {
	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		return nil, err
	}
	slog.Info("Database initialized successfully")

	err = db.AutoMigrate(&order.Order{}, &order.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	slog.Info("Database migrate successfully")

	return &App{
		DB:        db,
		OrderRepo: order.NewOrderRepository(db),
	}, nil
}

func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
