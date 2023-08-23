package routes

import (
	"fmt"
	"log"

	"example.com/m/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	go handlers.Gopython() // Ensure these are exported from handlers package
	go handlers.Gopythontwo()

	r.LoadHTMLGlob("/home/mik/Документы/GitHub/hrenovosti/templates/home.html")

	r.DELETE("/del/:id", handlers.DeleteArticleHandler)

	r.GET("/", handlers.RandomNewsHandler)
	r.GET("/articles", handlers.GetAllArticlesHandler)
	r.Static("/static", "./static")
	r.POST("/articles/add", handlers.AddArticleHandler)

	fmt.Println("Server started on http://localhost:80")
	log.Fatal(r.Run(":80"))

}
