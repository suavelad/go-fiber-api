package models

import "time"

type Product struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
