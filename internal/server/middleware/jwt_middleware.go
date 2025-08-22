package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/playground/userapi/pkg/utils"
)

// RequireAuth is an Echo middleware that validates the Authorization header (Bearer token),
// parses the JWT and sets the user id into the context under "user_id".
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header")
		}
		token := parts[1]
		uid, err := utils.ParseToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		// store in context for handlers
		c.Set("user_id", uid)
		return next(c)
	}
}

// Helper to retrieve user ID from context (returns 0 if absent).
func GetUserID(c echo.Context) int {
	v := c.Get("user_id")
	if v == nil {
		return 0
	}
	if id, ok := v.(int); ok {
		return id
	}
	return 0
}
