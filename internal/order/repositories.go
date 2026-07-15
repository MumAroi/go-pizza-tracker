package order

import (
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *Order) error
	GetOrder(id string) (*Order, error)
	GetOrders() ([]Order, error)
	UpdateOrderStatus(id string, status string) error
	DeleteOrder(id string) error
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

func (r *orderRepository) GetOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Preload("Items").Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateOrderStatus(id string, status string) error {
	err := r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
	return err
}

func (r *orderRepository) DeleteOrder(id string) error {
	err := r.db.Select("Items").Delete(&Order{ID: id}).Error
	return err
}
