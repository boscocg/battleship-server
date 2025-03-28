package main

import (
	"context"
	"log"

	config "battledak-server/configs"
	"battledak-server/internal/routes"
)

func main() {
	config.LoadEnv()
	ctx := context.Background()

	config.InitializeConfig(ctx)

	log.Printf("Environment: %s", config.GetEnv("ENV"))

	router := routes.SetupRouter()

	log.Fatal(router.Run(":8080"))
}
