package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/mygopher/gophertoken"
)

// AuthTokenFiberMiddleware checks for a role-based token and authorizes users
func AuthTokenFiberMiddleware(tokenManager gophertoken.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var role string
		var token string
		var found bool

		// Loop through possible roles to find the correct token
		for _, r := range configs.AllowedRoles {
			cookieName := r + "|_|" + configs.TokenBaseCookieName
			if cookieToken := c.Cookies(cookieName); cookieToken != "" {
				token = cookieToken
				role = r
				found = true
				break
			}
		}

		// If no token is found, return unauthorized
		if !found {
			response := responses.NewFiberResponse(
				c,
				fiber.StatusUnauthorized,
				"Unauthorized",
				nil,
				"Failed to find a valid role-based token in cookies",
			)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		// Validate the token
		payload, err := tokenManager.ValidateToken(token)
		if err != nil {
			response := responses.NewFiberResponse(
				c,
				fiber.StatusUnauthorized,
				"Invalid token",
				nil,
				err.Error(),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		// Store necessary information in the context
		c.Locals("payloadID", payload.ID)
		c.Locals("userID", payload.UserID)
		c.Locals("username", payload.Username)
		c.Locals("role", role) // Set the role in the context for authorization

		// Continue to the next middleware or handler
		return c.Next()
	}
}
