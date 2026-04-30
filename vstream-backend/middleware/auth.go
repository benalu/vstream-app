// vstream-backend/middleware/auth.go
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type adminClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// RequireAdmin memvalidasi JWT dari httpOnly cookie.
// Attach ke semua route /api/admin/* kecuali /login.
func RequireAdmin(c *gin.Context) {
	tokenStr, err := c.Cookie("vs_admin_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Sesi tidak ditemukan, silakan login",
		})
		return
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenStr, &adminClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Pastikan algoritma yang digunakan adalah HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Sesi tidak valid atau sudah expired",
		})
		return
	}

	claims, ok := token.Claims.(*adminClaims)
	if !ok || claims.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Akses ditolak",
		})
		return
	}

	// Simpan role ke context untuk dipakai handler jika perlu
	c.Set("role", claims.Role)
	c.Next()
}
