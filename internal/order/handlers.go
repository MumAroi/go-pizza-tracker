package order

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ServeNewOrderForm(c *gin.Context)
	HandleNewOrderPost(c *gin.Context)
	ServeCustomer(c *gin.Context)
}

type handler struct {
	order OrderRepository
}

func NewHandler(order OrderRepository) Handler {
	return &handler{order: order}
}

func (h *handler) ServeNewOrderForm(c *gin.Context) {
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes: GetPizzaTypes(),
		PizzaSizes: GetPizzaSizes(),
	})
}

func (h *handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderRequest

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItems := make([]OrderItem, len(form.Sizes))
	for i := range orderItems {
		orderItems[i] = OrderItem{
			Size:         form.Sizes[i],
			Pizza:        form.PizzaTypes[i],
			Instructions: form.Instructions[i],
			Quantity:     1,
		}
	}

	order := Order{
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Status:       GetOrderStatuses()[0],
		Items:        orderItems,
	}

	if err := h.order.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "error", err)
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}

	slog.Info("Order created", "orderId", order.ID, "customer", order.CustomerName)

	c.Redirect(http.StatusSeeOther, "/customer/"+order.ID)
}

func (h *handler) ServeCustomer(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := h.order.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "Order not found")
		return
	}

	c.HTML(http.StatusOK, "customer.tmpl", CustomerData{
		Title:    "Pizza Order Status " + orderID,
		Order:    *order,
		Statuses: GetOrderStatuses(),
	})

}
