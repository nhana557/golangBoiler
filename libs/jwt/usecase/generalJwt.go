package usecase

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SetJwtGeneral sets JWT validation for both admin and merchant access
func (h *JwtUsecase) SetJwtGeneral(g *gin.RouterGroup) {
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

	// Middleware to check access for admin or merchant
	g.Use(h.ValidateGeneralJwt)
}

// ValidateGeneralJwt checks if the token has admin access or is associated with a merchant
func (h *JwtUsecase) ValidateGeneralJwt(c *gin.Context) {
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

	// Check if the user is admin
	if isAdmin, ok := claims["is_admin"].(bool); ok && isAdmin {
		c.Next()
		return
	}

	// Otherwise, check for merchant ID (jti claim)
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

