package usecase

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SetJwtUser sets JWT validation middleware for user access only
func (h *JwtUsecase) SetJwtUser(g *gin.RouterGroup) {
	secret := h.Config.GetString("jwt.secret")

	// Middleware to validate JWT token
	g.Use(func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
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

	// Middleware to validate user ID in JWT claims
	g.Use(h.validateJwtUser)
}

// validateJwtUser is middleware that checks if the JWT contains a valid user ID (jti claim)
func (h *JwtUsecase) validateJwtUser(c *gin.Context) {
	userToken, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	token, ok := userToken.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token"})
		c.Abort()
		return
	}

	// Get user ID from the token's claims
	mid, ok := claims["jti"].(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token ID"})
		c.Abort()
		return
	}

	// Retrieve user information from the database
	ctx := context.TODO()
	user, err := h.getOneUser(ctx, mid)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}

	// Set user data in context if valid
	c.Set("user", user)
	c.Next()
}


