package session

import (
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"gorm.io/gorm"
)

func NewSessionStore(db *gorm.DB, secretKey []byte) sessions.Store {
	store := gormsessions.NewStore(db, true, secretKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: 3,
	})

	return store
}
