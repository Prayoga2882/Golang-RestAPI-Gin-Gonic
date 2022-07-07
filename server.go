package main

import (
	"main/config"
	"main/middleware"
	"main/routes"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main(){
	config.InitDB()
	
	gotenv.Load()
	router := gin.Default()
	
	// Routes
	router.POST("/Register", routes.RegisterHandler)
	router.GET("/Login", routes.LoginHandler)

	router.GET("/Check", middleware.IsAuth(), routes.CheckToken)

	router.GET("/auth/:provider", routes.RedirectHandler)
	router.GET("/auth/:provider/callback", routes.CallbackHandler)
	
	router.GET("/", routes.Base)
	router.GET("/GetArticles", routes.GetArticles)
	router.GET("/GetArticle/:id", routes.GetArticle)
	router.POST("/PostArticle/", routes.PostArticle)
	router.PUT("/UpdateArticle/:id", middleware.IsAuth(), routes.UpdateArticle)
	router.DELETE("/DeleteArticle/:id", middleware.IsAuth(), routes.DeleteArticle)

	router.Run(":8080")
}