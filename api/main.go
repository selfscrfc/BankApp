package main

import (
	"PetBank/config"
	"PetBank/pkg/logger"
	"PetBank/pkg/routes"
	"flag"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic("Parse config " + err.Error())
	}

	log, err := logger.SetupLogger(cfg.Server.Mode)

	if err != nil {
		panic("Parse logger " + err.Error())
	}

	log.Infof("App started on port %s. %s %s",
		cfg.GRPC.Port,
		cfg.Server.Mode,
		cfg.Server.AppVersion)

	serverAddr := flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	customerService, err := grpc.Dial(*serverAddr)

	app := fiber.New()

	routes.PrivateRoutes(app, customerService)
	routes.PublicRoutes(app, customerService)

	app.Listen(":3000")
}
