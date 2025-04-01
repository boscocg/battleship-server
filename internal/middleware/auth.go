package middleware

import (
	config "battledak-server/configs"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			allowedAPIKeys := strings.Split(config.GetEnv("ALLOWED_API_KEYS"), ",")
			for _, key := range allowedAPIKeys {
				if apiKey == strings.TrimSpace(key) {
					c.Next()
					return
				}
			}
		}

		// allowedIPs := config.GetEnv("ALLOWED_IPS")
		// if allowedIPs != "" {
		// 	clientIP := c.ClientIP()
		// 	allowedIPList := strings.Split(allowedIPs, ",")

		// 	for _, ipRange := range allowedIPList {
		// 		ipRange = strings.TrimSpace(ipRange)

		// 		if strings.Contains(ipRange, "/") {
		// 			_, ipNet, err := net.ParseCIDR(ipRange)
		// 			if err == nil {
		// 				clientIPAddr := net.ParseIP(clientIP)
		// 				if ipNet.Contains(clientIPAddr) {
		// 					c.Next()
		// 					return
		// 				}
		// 			}
		// 		} else if clientIP == ipRange {
		// 			c.Next()
		// 			return
		// 		}
		// 	}
		// }

		// allowedOrigins := config.GetEnv("ALLOWED_ORIGINS")
		// if allowedOrigins != "" {
		// 	origin := c.GetHeader("Origin")
		// 	if origin != "" {
		// 		originList := strings.Split(allowedOrigins, ",")
		// 		for _, allowedOrigin := range originList {
		// 			allowedOrigin = strings.TrimSpace(allowedOrigin)
		// 			if strings.HasSuffix(origin, allowedOrigin) {
		// 				c.Next()
		// 				return
		// 			}
		// 		}
		// 	}
		// }

		// if config.GetEnv("ENV") == "dev" && config.GetEnv("DEV_AUTH_BYPASS") == "true" {
		// 	c.Next()
		// 	return
		// }

		c.AbortWithStatusJSON(403, gin.H{
			"code":    "ACCESS_DENIED",
			"message": "You are not authorized to access this API",
		})
	}
}

func PublicAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
