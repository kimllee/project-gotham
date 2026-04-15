package routes

import (
	"project-gotham/controllers"
	"project-gotham/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	routesGroup := router.Group("/comments")
	routesGroup.Use(middlewares.Authentication()) // On ne peut plus utiliser les routes suivantes sans être authentifié.

	{
		routesGroup.POST("/", controllers.PostComment)
	}
}
