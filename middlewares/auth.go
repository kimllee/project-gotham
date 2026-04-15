package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Accès non autorisé."})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Structure "t" qui renvoie une interface
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// On attend une signature HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			// Retour en byte de la signature générée précédemment.
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Jeton invalide ou expiré."})
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Impossible de lire le jeton."})
			return
		}

		userID := int(claim["UserID"].(float64)) // Transformation en int le UserID
		c.Set("userID", userID)                  // On l'ajoute à la clé userID
		c.Next()                                 // On passe à la suite.
	}
}
