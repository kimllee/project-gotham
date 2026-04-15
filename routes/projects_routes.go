package routes

import (
	"project-gotham/controllers"
	"project-gotham/middlewares"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(router *gin.Engine) {
	routesGroup := router.Group("/projects")

	routesGroup.Use(middlewares.Authentication()) // On ne peut plus utiliser les routes suivantes sans être authentifié.

	{
		// Lecture séquentielle c'est-à-dire que si la route LikedProjects est après la route PutProjects elle ne sera jamais évaluée.
		routesGroup.GET("/", controllers.GetProjects)
		routesGroup.POST("/", controllers.PostProject)
		routesGroup.GET("/:id", controllers.GetProject)
		routesGroup.PUT("/:id/like", controllers.LikedProjects)
		routesGroup.PUT("/:id", controllers.PutProject)
		routesGroup.DELETE("/:id", controllers.DeleteProject)

	}

}
