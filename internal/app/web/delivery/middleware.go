package delivery

import (
	"errors"
	"go-Chatting/utils"
	"net/http"
	"time"

	webdomain "go-Chatting/domain/app/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey *utils.SecretKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// Verify the token
		userID, err := verifyToken(cookie, []byte(secretKey.Get()))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// Set the user ID in the context for future use
		c.Set("userID", userID)

		// Call the next middleware/handler
		c.Next()
	}
}

func verifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
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

	return token, nil
}
