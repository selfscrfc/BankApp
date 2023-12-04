package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/selfscrfc/PetBank/api/docs"
	"github.com/selfscrfc/PetBank/config"
	mygrpc "github.com/selfscrfc/PetBank/internal/grpc"
	"github.com/selfscrfc/PetBank/pkg/logger"
	"github.com/selfscrfc/PetBank/pkg/routes"
)

// @title PetBank API
// @version 1.0
// @license.name Apache 2.0
// @PetBank service
// @host localhost:3000
// @BasePath /
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

	app := fiber.New()

	customerClient, err := mygrpc.NewCustomerClient(cfg)

	if err != nil {
		log.Error("Customer client connection error " + err.Error())
	}

	accountsClient, err := mygrpc.NewAccountsClient(cfg)

	if err != nil {
		log.Error("Accounts client connection error " + err.Error())
	}

	routes.PrivateRoutes(app, customerClient, accountsClient)
	routes.PublicRoutes(app, customerClient, accountsClient)

	app.Get("/swagger/*", swagger.HandlerDefault)
	if err = app.Listen(":3000"); err != nil {
		log.Error("Server listen error ", err.Error())
	}
}
