package models

import "time"

type Order struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ProductRefer int       `json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductRefer"`
	UserRefer    int       `json:"user_id"`
	User         User      `gorm:"foreignKey:UserRefer"`
}
