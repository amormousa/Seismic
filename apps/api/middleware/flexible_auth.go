package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/models"
)

// RequireAuthOrAPIKey accepts either a JWT access token or an
// API key in the Authorization header. Used by endpoints that
// both the web dashboard and editor plugins need to call.
func RequireAuthOrAPIKey(pool *pgxpool.Pool, jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return helpers.Error(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		value := strings.TrimPrefix(header, "Bearer ")
		ctx := c.Context()

		// Try as JWT first
		token, err := jwt.Parse(value, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})
		if err == nil && token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				if userID, ok := claims["sub"].(string); ok {
					c.Locals("userID", userID)
					return c.Next()
				}
			}
		}

		// Not a valid JWT, try as API key
		user, err := models.FindUserByAPIKey(ctx, pool, value)
		if err != nil {
			return helpers.Error(c, fiber.StatusInternalServerError, "Something went wrong")
		}
		if user == nil {
			return helpers.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
		}

		c.Locals("userID", user.ID)
		return c.Next()
	}
}
