package models

import "github.com/jinzhu/gorm"

type Articles struct {
	gorm.Model
	Title 		string `json:"title" binding:"required"`
	Price 		int		`json:"price" binding:"required"`
	Description string 	`json:"description" binding:"required"`
}