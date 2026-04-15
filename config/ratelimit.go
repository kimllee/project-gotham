package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimit(reqpersec int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(reqpersec), reqpersec) // nbr req, nbr max req

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Trop de requêtes lancées."})
			return
		}
		c.Next()
	}
}
