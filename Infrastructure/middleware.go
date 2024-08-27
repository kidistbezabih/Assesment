package infrastructure

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "auth header"})
			c.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unautorized access"})
			c.Abort()
			return
		}
		tokenClaim := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(authParts[1], tokenClaim, validateSigningMethod)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("user_id", tokenClaim["id"])
		c.Set("is_admin", tokenClaim["isadmin"])
		c.Next()
	}
}

func AdminMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// uses_name, user_exist := c.Get("userame")
		role, role_exist := c.Get("is_admin")
		roleParsed := role.(bool)

		if !role_exist || !roleParsed {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unautorized access"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 		c.Set("loan_id", tokenClaim["loanid"])

func validateSigningMethod(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("SECRET_KEY")), nil
}
