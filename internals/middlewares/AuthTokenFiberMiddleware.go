package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func AuthTokenFiberMiddleware(tokenManager gophertoken.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from the cookie
		token := c.Cookies("SuperUserAuthorizationToken")
		if token == "" {
			response := responses.NewFiberResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil, "Failed to get token from cookie")
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		// Validate the token
		payload, err := tokenManager.ValidateToken(token)
		if err != nil {
			response := responses.NewFiberResponse(c, fiber.StatusUnauthorized, "Invalid token", nil, err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		// Store necessary information in the context
		c.Locals("userID", payload.ID)
		c.Locals("username", payload.Username)

		// Continue to the next middleware or handler
		return c.Next()
	}
}
