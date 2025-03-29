package main

import (
	"context"
	"log"

	config "battledak-server/configs"
	"battledak-server/internal/routes"
)

func main() {
	config.LoadEnv()
	log.Printf("Environment: %s", config.GetEnv("ENV"))

	ctx := context.Background()

	config.InitializeConfig(ctx)

	router := routes.SetupRouter()

	log.Fatal(router.Run(":8080"))
}
