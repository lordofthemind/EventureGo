package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/responses"
	"github.com/lordofthemind/mygopher/gophertoken"
)

// AuthTokenGinMiddleware checks for a role-based token and authorizes users
func AuthTokenGinMiddleware(tokenManager gophertoken.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role string
		var token string
		var found bool

		// Loop through possible roles to find the correct token
		for _, r := range configs.AllowedRoles {
			cookieName := r + "|_|" + configs.TokenBaseCookieName
			if cookieToken, err := c.Cookie(cookieName); err == nil {
				token = cookieToken
				role = r
				found = true
				break
			}
		}

		// If no token is found, return unauthorized
		if !found {
			response := responses.NewGinResponse(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				nil,
				"Failed to find a valid role-based token in cookies",
			)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Validate the token
		payload, err := tokenManager.ValidateToken(token)
		if err != nil {
			response := responses.NewGinResponse(
				c,
				http.StatusUnauthorized,
				"Invalid token",
				nil,
				err.Error(),
			)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Set values in the context for use in the rest of the application
		c.Set("payloadID", payload.ID)
		c.Set("userID", payload.UserID)
		c.Set("username", payload.Username)
		c.Set("role", role)

		// Proceed with the request
		c.Next()
	}
}
