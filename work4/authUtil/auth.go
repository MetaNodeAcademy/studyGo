package authUtil

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"work4/model"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header获取token，如果没有则从cookie获取
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" {
			c.JSON(401, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		claims, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !claims.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claimsMap, ok := claims.Claims.(jwt.MapClaims); ok {
			userID := uint(claimsMap["id"].(float64))
			username := claimsMap["username"].(string)

			// 验证用户ID在数据库中是否存在
			if !validateUserExists(userID) {
				c.JSON(401, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			c.Set("userID", userID)
			c.Set("username", username)
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}

// validateUserExists 验证用户ID在数据库中是否存在
func validateUserExists(userID uint) bool {
	dbManager := model.GetDBManager()
	var user model.User
	err := dbManager.DB.First(&user, userID).Error
	return err == nil
}
