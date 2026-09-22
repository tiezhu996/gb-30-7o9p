package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/util"
)

// RequestIDKey is the gin context key for the request id.
const RequestIDKey = "request_id"

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}

// RequestLogger logs structured request metadata.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := newRequestID()
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
		rl := util.LoggerWithRequest(logger, requestID, c.Request.Method, c.Request.URL.Path)
		rl.Info(fmt.Sprintf(constants.LogRequestHandled, requestID, c.Request.Method, c.Request.URL.Path,
			c.Writer.Status(), time.Since(start).Milliseconds()))
	}
}
