package orders

import (
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetOrder(id string) (*Order, error) {
	var order Order
	err := r.db.Preload("Items").First(&order, "id = ?", id).Error
	return &order, err
}
