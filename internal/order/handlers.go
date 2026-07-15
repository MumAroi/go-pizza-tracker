package order

import (
	"log/slog"
	"net/http"
	"pizza-tracker/internal/shared/notification"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ServeNewOrderForm(c *gin.Context)
	HandleNewOrderPost(c *gin.Context)
	ServeInfo(c *gin.Context)
	GetNotification(c *gin.Context)
}

type handler struct {
	OrderDeps
}

type OrderDeps struct {
	OrderRepo       OrderRepository
	NotificationMgr *notification.NotificationManager
}

func NewHandler(deps OrderDeps) Handler {
	return &handler{
		OrderDeps: deps,
	}
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

	if err := h.OrderRepo.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "error", err)
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}

	slog.Info("Order created", "orderId", order.ID, "customer", order.CustomerName)

	h.NotificationMgr.Notify("admin:new_orders", "new_orders")

	c.Redirect(http.StatusSeeOther, "/orders/"+order.ID)
}

func (h *handler) ServeInfo(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := h.OrderRepo.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "Order not found")
		return
	}

	c.HTML(http.StatusOK, "info.tmpl", OrderInfoData{
		Title:    "Pizza Order Status " + orderID,
		Order:    *order,
		Statuses: GetOrderStatuses(),
	})

}

func (h *handler) GetNotification(c *gin.Context) {
	orderID := c.Query("orderId")

	if orderID == "" {
		c.String(400, "Invalid orderId")
		return
	}

	_, err := h.OrderRepo.GetOrder(orderID)
	if err != nil {
		c.String(404, "Order not found")
		return
	}

	key := "order:" + orderID
	client := make(chan string, 10)

	h.NotificationMgr.AddClient(key, client)

	defer func() {
		h.NotificationMgr.RemoveClient(key, client)
		slog.Info("Customer client disconnected", "orderId", orderID)
	}()

	h.NotificationMgr.StreamSSE(c, client)
}
