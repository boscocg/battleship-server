package routes

import (
	"battledak-server/internal/controller"
	"battledak-server/internal/middleware"
	"battledak-server/internal/service"
	"strings"
	"time"

	config "battledak-server/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if config.GetEnv("ENV") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	trustedProxies := config.GetEnv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxyList := strings.Split(trustedProxies, ",")
		for i := range proxyList {
			proxyList[i] = strings.TrimSpace(proxyList[i])
		}
		router.SetTrustedProxies(proxyList)
	} else if config.GetEnv("ENV") == "dev" {
		router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	} else {
		// In production, we don't set any trusted proxies
		router.SetTrustedProxies([]string{})
	}

	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if config.GetEnv("ENV") == "dev" || config.GetEnv("ENV") == "local" {
				return true
			}
			if strings.HasSuffix(origin, "-project.vercel.app") || strings.HasSuffix(origin, ".project.tech") {
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
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cryptoService := service.NewCryptoService()
	gameService := service.NewGameService()
	moveService := service.NewMoveService()
	gameController := controller.NewGameController(gameService)
	moveController := controller.NewMoveController(moveService, gameService, cryptoService)
	healthController := controller.NewHealthController()

	public := router.Group("/")
	public.Use(middleware.PublicAccess())
	{
		public.GET("/health", healthController.Check)
	}

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		game := protected.Group("/game")
		{
			game.GET("/:id", gameController.GetGame)
			game.POST("", gameController.StartGame)
		}
		move := protected.Group("/move")
		{
			move.POST("", moveController.Move)
		}
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
