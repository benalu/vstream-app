// vstream-backend/handlers/auth.go
package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type adminClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Login menerima admin key, return JWT via httpOnly cookie
func Login(c *gin.Context) {
	var input struct {
		Key string `json:"key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Key tidak boleh kosong",
		})
		return
	}

	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Server belum dikonfigurasi",
		})
		return
	}

	// Validasi key (constant-time compare untuk mencegah timing attack)
	if !secureCompare(input.Key, adminKey) {
		// Delay kecil untuk mencegah brute-force
		time.Sleep(500 * time.Millisecond)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Admin key salah",
		})
		return
	}

	// Buat JWT token
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := adminClaims{
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vstream",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal membuat session",
		})
		return
	}

	// Set sebagai httpOnly cookie
	// SameSite=Lax cukup untuk SPA same-origin; ubah ke None + Secure jika beda domain
	c.SetCookie(
		"vs_admin_token", // name
		tokenStr,         // value
		60*60*24,         // maxAge: 24 jam (detik)
		"/",              // path
		"",               // domain (kosong = current domain)
		false,            // secure (set true jika HTTPS)
		true,             // httpOnly — tidak bisa diakses JS
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
	})
}

// Logout menghapus cookie
func Logout(c *gin.Context) {
	c.SetCookie("vs_admin_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}

// Me mengecek apakah session masih valid (dipakai frontend untuk route guard)
func Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"role": "admin",
		},
	})
}

// secureCompare membandingkan dua string dengan waktu konstan
// untuk mencegah timing attack
func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
