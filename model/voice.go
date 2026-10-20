package model

type Voice struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:100;not null"`
	Gender      string `gorm:"size:20;not null"`
	IsAvailable bool   `gorm:"not null;default:false"`
}
