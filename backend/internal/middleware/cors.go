package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/config"
)

// CORS returns a CORS middleware whose allowed origins come from
// APP_CORS_ORIGINS (comma-separated). The wildcard "*" is only enabled when
// explicitly configured; otherwise the configured origin list is used.
func CORS(cfg *config.Config) gin.HandlerFunc {
	raw := strings.TrimSpace(cfg.CORSOrigins)
	if raw == "" {
		raw = "http://localhost:8101"
	}

	if raw == "*" {
		return cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"X-Request-Id"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}

	origins := make([]string, 0)
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:8101"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
