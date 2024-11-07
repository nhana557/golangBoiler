package usecase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// SetJwtAdmin sets middleware for JWT validation and admin-only access
func (h *JwtUsecase) SetJwtAdmin(g *gin.RouterGroup) {
	secret := h.Config.GetString("jwt.secret")

	// Middleware for JWT token validation
	g.Use(func(c *gin.Context) {
		// Extract and validate JWT token
		tokenString := c.GetHeader("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set the token in context if it's valid
		c.Set("user", token)
		c.Next()
	})

	// Middleware to validate admin access
	g.Use(h.validateJwtAdmin)
}

// validateJwtAdmin is middleware that allows only admin access
func (h *JwtUsecase) validateJwtAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !claims.VerifyIssuer("admin", true) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	isAdmin, ok := claims["is_admin"].(bool)
	if !ok || !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	// Continue to the next middleware or handler if the user is an admin
	c.Next()
}
