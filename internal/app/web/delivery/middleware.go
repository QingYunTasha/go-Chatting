package delivery

import (
	"errors"
	"net/http"
	"time"

	webdomain "go-Chatting/domain/app/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey webdomain.SecretKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("id")
		cookie, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// Verify the token
		token, err := verifyToken(cookie, uint32(userID), []byte(secretKey.Get()))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// Set the user ID in the context for future use
		c.Set("token", token)

		// Call the next middleware/handler
		c.Next()
	}
}

func verifyToken(tokenString string, userID uint32, secretKey []byte) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &webdomain.CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check the token signature
	if !token.Valid {
		return nil, errors.New("invalid token signature")
	}

	// Check the token expiry time
	claims, ok := token.Claims.(*webdomain.CustomJWTClaims)
	if !ok || claims.ExpiresAt.Unix() < time.Now().UTC().Unix() {
		return nil, errors.New("token has expired")
	}

	if claims.UserID != userID {
		return nil, errors.New("invalid authentication")
	}

	return token, nil
}
