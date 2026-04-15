package config

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func SecurityMiddleware() gin.HandlerFunc {
	sec := secure.New(secure.Options{
		FrameDeny:             true,  // Click-jacking
		ContentTypeNosniff:    true,  // MIME sniffing - tentatives d'envoi d'un contenu en apparence JSON mais qui est réellement du code exécutable malveillant.
		BrowserXssFilter:      true,  // Attaques XSS : Cross Site Scripping
		SSLRedirect:           false, // Forçage du https
		ContentSecurityPolicy: "default-src 'self'",
	})

	return func(c *gin.Context) {
		err := sec.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}
