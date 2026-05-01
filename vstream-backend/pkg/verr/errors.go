package verr

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Status      int
	UserMessage string // pesan aman untuk user / admin UI
	TechDetail  string // detail teknis dari original error
	LogMsg      string // untuk slog
	IsAdmin     bool
}

// NewAdminError — route /api/admin/*
// Admin dapat pesan + detail teknis dari err.Error()
func NewAdminError(status int, userMessage string, err error) *AppError {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	return &AppError{
		Status:      status,
		UserMessage: userMessage,
		TechDetail:  detail,
		LogMsg:      userMessage,
		IsAdmin:     true,
	}
}

// NewPublicError — route publik
// Frontend hanya dapat publicMsg, detail teknis hanya masuk log
func NewPublicError(status int, publicMsg string, logMsg string, err error) *AppError {
	techDetail := ""
	if err != nil {
		techDetail = err.Error()
	}
	return &AppError{
		Status:      status,
		UserMessage: publicMsg,
		TechDetail:  techDetail,
		LogMsg:      logMsg,
		IsAdmin:     false,
	}
}

func Handle(c *gin.Context, appErr *AppError) {
	// Log selalu — termasuk detail teknis
	attrs := []any{
		"status", appErr.Status,
		"msg", appErr.LogMsg,
	}
	if appErr.TechDetail != "" {
		attrs = append(attrs, "error", appErr.TechDetail)
	}
	slog.Error("app_error", attrs...)

	// Admin: dapat pesan + detail teknis (kalau ada)
	if appErr.IsAdmin {
		resp := gin.H{
			"success": false,
			"error":   appErr.UserMessage,
		}
		if appErr.TechDetail != "" {
			resp["detail"] = appErr.TechDetail
		}
		c.JSON(appErr.Status, resp)
		return
	}

	// Public: hanya pesan aman, tidak ada detail
	c.JSON(appErr.Status, gin.H{
		"success": false,
		"error":   appErr.UserMessage,
	})
}
