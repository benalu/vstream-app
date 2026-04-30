package verr

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Status  int    `json:"-"`
	Message string `json:"message"` // Pesan aman untuk user
	LogMsg  string `json:"-"`       // Detail teknis untuk slog
	Err     error  `json:"-"`
}

func Handle(c *gin.Context, err *AppError) {
	slog.Error(err.LogMsg, "status", err.Status, "error", err.Err)
	c.JSON(err.Status, gin.H{"success": false, "error": err.Message})
}

// Error khusus Admin (Bisa lebih detail)
func NewAdminError(status int, logMsg string, err error) *AppError {
	return &AppError{Status: status, Message: "Admin Action Failed: " + logMsg, LogMsg: logMsg, Err: err}
}

// Error khusus Frontend (Generic/Aman)
func NewPublicError(status int, publicMsg string, logMsg string, err error) *AppError {
	return &AppError{Status: status, Message: publicMsg, LogMsg: logMsg, Err: err}
}
