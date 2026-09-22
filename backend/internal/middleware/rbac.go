package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
)

// RequireRole rejects requests whose role is not allowed.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		if !allowed[GetUserRole(c)] {
			c.AbortWithStatusJSON(http.StatusForbidden,
				dto.Fail(constants.CodeForbidden, constants.MsgForbidden+": require role "+roles[0]))
			return
		}
		c.Next()
	}
}
