package notification

import (
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

type NotificationManager struct {
	clients map[string]map[chan string]bool
	mu      sync.RWMutex
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		clients: make(map[string]map[chan string]bool),
	}
}

func (n *NotificationManager) AddClient(key string, client chan string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.clients[key] == nil {
		n.clients[key] = make(map[chan string]bool)
	}

	n.clients[key][client] = true
}

func (n *NotificationManager) RemoveClient(key string, client chan string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if clients := n.clients[key]; clients != nil {
		delete(clients, client)
		if len(clients) == 0 {
			delete(n.clients, key)
		}
	}

	close(client)
}

func (n *NotificationManager) Notify(key, message string) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for client := range n.clients[key] {
		select {
		case client <- message:
		default:
		}
	}
}

func (n *NotificationManager) StreamSSE(c *gin.Context, client chan string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-client; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
