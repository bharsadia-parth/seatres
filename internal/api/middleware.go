package api

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID uses the caller's X-Request-ID or generates one,
// stores it in the context, and echoes it back in the response.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// AccessLog writes one structured JSON log line per request.
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.Info("request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}

// Auth reads "Authorization: Bearer <token>". The token IS the identity.
// A token equal to adminToken also gets admin rights.
// user_id is NEVER read from the request body.
func Auth(adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		if !strings.HasPrefix(h, "Bearer ") || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_or_invalid_token"})
			return
		}
		c.Set("user_id", token)
		c.Set("is_admin", adminToken != "" && token == adminToken)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("is_admin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin_only"})
			return
		}
		c.Next()
	}
}
