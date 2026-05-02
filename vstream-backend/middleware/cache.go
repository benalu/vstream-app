package middleware

import "github.com/gin-gonic/gin"

// PublicCache menambahkan header Cache-Control untuk public endpoints.
// max-age dalam detik.
func PublicCache(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=30")
		c.Next()
	}
}
