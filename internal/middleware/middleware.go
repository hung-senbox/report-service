package middleware

import (
	"context"
	"net/http"
	"report-service/helper"
	"report-service/pkg/constants"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Secured() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		// app language header
		appLanguage := helper.ParseAppLanguage(c.GetHeader("X-App-Language"), 1)
		c.Writer.Header().Set("X-App-Language", strconv.Itoa(int(appLanguage)))

		// context gốc
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constants.AppLanguage, appLanguage)
		c.Set(constants.AppLanguage.String(), appLanguage)

		if len(authorizationHeader) == 0 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims[constants.UserID.String()].(string); ok {
				c.Set(constants.UserID.String(), userId)
				ctx = context.WithValue(ctx, constants.UserID, userId)
			}
			if userName, ok := claims[constants.UserName.String()].(string); ok {
				c.Set(constants.UserName.String(), userName)
				ctx = context.WithValue(ctx, constants.UserName, userName)
			}
			if userRoles, ok := claims[constants.UserRoles.String()].(string); ok {
				c.Set(constants.UserRoles.String(), userRoles)
				ctx = context.WithValue(ctx, constants.UserRoles, userRoles)
			}
		}

		// Token
		c.Set(constants.Token.String(), tokenString)
		ctx = context.WithValue(ctx, constants.Token, tokenString)

		// Gán lại context duy nhất cho request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesAny, exists := c.Get(constants.UserRoles.String())
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Roles not found"})
			return
		}

		rolesStr, ok := rolesAny.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid roles format"})
			return
		}

		// ví dụ roles: "SuperAdmin, Teacher"
		roles := strings.Split(rolesStr, ",")
		isAdmin := false
		for _, role := range roles {
			if strings.TrimSpace(role) == "SuperAdmin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}
