package routes

import (
	"battledak-server/internal/controller"
	"strings"
	"time"

	config "battledak-server/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if config.GetEnv("ENV") == "dev" {
				return true
			}
			if strings.HasSuffix(origin, "-gateway-dao.vercel.app") || strings.HasSuffix(origin, ".gateway.tech") {
				return true
			}
			allowedOrigins := []string{config.GetEnv("LOCAL_FRONT"), config.GetEnv("API_HOST_FULL")}
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "OPTIONS", "PATCH", "PUT", "DELETE", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	gameController := controller.NewGameController()

	game := router.Group("/game")
	{
		game.GET("/:id", gameController.GetGame)
		game.POST("", gameController.StartGame)
	}

	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.JSON(404, gin.H{"code": "ENDPOINT_NOT_FOUND", "message": "Endpoint not found"})
	})

	return router
}
