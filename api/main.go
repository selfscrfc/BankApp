package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/selfscrfc/PetBank/config"
	_ "github.com/selfscrfc/PetBank/docs"
	mygrpc "github.com/selfscrfc/PetBank/internal/grpc"
	"github.com/selfscrfc/PetBank/pkg/logger"
	"github.com/selfscrfc/PetBank/pkg/routes"
	Customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
	log "log"
)

// @title PetBank API
// @version 1.0
// @license.name Apache 2.0
// @PetBank service
// @host localhost:3000
// @BasePath /
func main() {
	log.Println("Starting API server")

	cfg, err := config.LoadConfig()

	if err != nil {
		panic("Parse config " + err.Error())
	}

	log, err := logger.SetupLogger(cfg.Server.Mode)

	if err != nil {
		panic("Parse logger " + err.Error())
	}

	app := fiber.New()

	customerClient, err := mygrpc.NewCustomerClient(cfg)

	_, err = (*customerClient).Create(context.Background(), &Customers.CreateRequest{
		FullName: "123",
		Login:    "123",
		Password: "123",
	})

	if err != nil {
		log.Error("Customer client connection error " + err.Error())
	}

	accountsClient, err := mygrpc.NewAccountsClient(cfg)

	if err != nil {
		log.Error("Accounts client connection error " + err.Error())
	}

	routes.PrivateRoutes(app, customerClient, accountsClient)
	routes.PublicRoutes(app, customerClient, accountsClient)

	log.Infof("App info: port :%s Mode: %s Version: %s",
		cfg.Server.Port,
		cfg.Server.Mode,
		cfg.Server.AppVersion)

	app.Get("/swagger/*", swagger.HandlerDefault)
	if err = app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Error("Server listen error ", err.Error())
	}

}
