package models

import (
	"errors"
	"time"
)

type User struct {
    ID        uint64 `json:"id" gorm:"primary_key"`
    UpdatedAt time.Time `json:"-"`
    // DeletedAt *time.Time `json:"-"`

    Login string `json:"login" gorm:"not null"`
    Name  string `json:"name" gorm:"not null"`

    // Wallets []Wallet `json:"wallets"`
}

func (u *User) Create() error {
	if err := u.validateCreation(); err != nil {
		return err
	}
	return database.Create(&u).Error
}

func (u *User) GetByLogin() bool {
	query := database.Where("login = ?", u.Login).First(&u)
	return query.Error == nil && !query.RecordNotFound()
}

func (u *User) validateCreation() error {
	if u.Login == "" {
		return errors.New("The field 'login' from User must have a value")
	}
	if u.Name == "" {
		return errors.New("The field 'name' from User must have a value")
	}
	return nil
}