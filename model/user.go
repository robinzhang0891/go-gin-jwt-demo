package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null"`
	Password string `gorm:"size:100;not null"`
}
