package config

import (
	"main/models"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	dsn := "root:@tcp(127.0.0.1:3306)/gin-api?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed connect to database")
	}
	
	DB.AutoMigrate(&models.Login{})
	DB.AutoMigrate(&models.Articles{})
	
}