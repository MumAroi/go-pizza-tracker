package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, username, password string) error {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := User{Username: username, Password: string(hashed)}
	return db.Create(&admin).Error
}
