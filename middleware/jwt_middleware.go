package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anhhuy1010/DATN-cms-ideas/config"
	"github.com/anhhuy1010/DATN-cms-ideas/helpers/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		tokenStr := c.GetHeader("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		cfg := config.GetConfig()
		jwtKey := []byte(cfg.GetString("auth.key"))

		claims := &util.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ✅ Token hợp lệ, lưu thông tin vào context
		c.Set("customer_uuid", claims.Uuid)
		c.Set("customer_name", claims.UserName)
		c.Set("customer_email", claims.Email)
		c.Set("start_day", claims.StartDay)
		c.Set("end_day", claims.EndDay)
		c.Next()
	}
}

func JWTDateCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startDayVal, startDayExists := c.Get("start_day")
		endDayVal, endDayExists := c.Get("end_day")
		if !startDayExists || !endDayExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing startDay or endDay"})
			c.Abort()
			return
		}

		startDay, okStart := startDayVal.(*time.Time)
		endDay, okEnd := endDayVal.(*time.Time)
		if !okStart || !okEnd || startDay == nil || endDay == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid startDay or endDay"})
			c.Abort()
			return
		}

		now := time.Now()
		if now.Before(*startDay) || now.After(*endDay) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or not yet valid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
