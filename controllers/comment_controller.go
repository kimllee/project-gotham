package controllers

import (
	"net/http"
	"project-gotham/config"
	"project-gotham/models"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	// Récupération de l'ID de l'utilisateur
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'ID de l'utilisateur est introuvable."})
		return
	}

	userIDInt, ok := userID.(int) // Conversion en int
	comment.UserID = uint(userIDInt)

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de l'enregistrement du commentaire."})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
