package controllers

import (
	"net/http"
	"os"
	"project-gotham/config"
	"project-gotham/models"
	"project-gotham/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Les données saisies sont invalides."})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'email ou le mot de passe est incorrect."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'email ou le mot de passe est incorrect."})
		return
	}

	claim := &CustomClaim{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de la génération du token."})
		return
	}

	c.JSON(http.StatusOK, tokenString)
}

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email déjà utilisé."})
		return
	}

	if err := utils.ValidatePassword(user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne serveur."})
		return
	}

	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'utilisateur."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur créé avec succès."})

}
