package main

import (
	"fmt"
	"log"
	"project-gotham/config"
	"project-gotham/models"
	"project-gotham/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "project-gotham/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Project Gotham
// @version 1.0
// @description Description du projet project Gotham
// @securityDefinition.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))

	err := godotenv.Load()
	if err != nil {
		log.Println("Aucun fichier .env trouvé.")
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong-pong"})
	})

	routes.ProjectRoutes(router)
	routes.UserRoutes(router)
	routes.CommentRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config.ConnectDB()
	fmt.Println("🏎️  Serveur démarré sur http://localhost:8000")
	//	http.ListenAndServe(":8080", nil)

	// Sert à créer une table SQL "projects" au démarrage. Gorm comprend quelle table utiliser en fonction du nom de la structure définie appelée.
	// Au pluriel sans miniscule, c'est une norme utilisée par Gorm.
	config.DB.AutoMigrate(&models.Project{}, &models.User{}, &models.Comment{})

	router.Run(":8000")

}
