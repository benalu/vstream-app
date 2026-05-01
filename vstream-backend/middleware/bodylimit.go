package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit membatasi ukuran request body.
// maxBytes: misal 2<<20 = 2MB, 10<<20 = 10MB
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
