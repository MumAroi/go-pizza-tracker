package app

import (
	"fmt"
	"log/slog"
	"pizza-tracker/internal/database"
	"pizza-tracker/internal/order"
	"pizza-tracker/internal/shared/notification"
	"pizza-tracker/internal/user"

	"gorm.io/gorm"
)

type App struct {
	DB              *gorm.DB
	OrderRepo       order.OrderRepository
	UserRepo        user.UserRepository
	NotificationMgr *notification.NotificationManager
}

func NewApp(dbPath string) (*App, error) {
	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		return nil, err
	}
	slog.Info("Database initialized successfully")

	err = db.AutoMigrate(&order.Order{}, &order.OrderItem{}, &user.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	slog.Info("Database migrate successfully")

	user.SeedAdmin(db, "admin", "Pas332211")
	slog.Info("Admin user seeded successfully")

	return &App{
		DB:              db,
		OrderRepo:       order.NewOrderRepository(db),
		UserRepo:        user.NewUserRepository(db),
		NotificationMgr: notification.NewNotificationManager(),
	}, nil
}

func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
