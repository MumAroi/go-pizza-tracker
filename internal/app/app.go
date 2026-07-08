package app

import (
	"fmt"
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

	err = db.AutoMigrate(&order.Order{}, &order.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

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
