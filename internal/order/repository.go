package order

import (
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *Order) error
	GetOrder(id string) (*Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetOrder(id string) (*Order, error) {
	var order Order
	err := r.db.Preload("Items").First(&order, "id = ?", id).Error
	return &order, err
}
