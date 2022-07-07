package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Articles []Articles 
	FullName 	string	
	Username  	string
	Password 	string 	
	Email 		string 
	Social_Id 	string
	Provider 	string
	Avatar 		string
}

type Login struct {
	
	Id       int64  `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}