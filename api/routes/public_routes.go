package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/selfscrfc/PetBank/api/controllers"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

func PublicRoutes(a *fiber.App, customerClient *customers.CustomerClient,
	accountsClient *account.AccountServiceClient) {

	a.Post("/signup", func(ctx *fiber.Ctx) error {
		return controllers.CreateCustomer(ctx, customerClient)
	})
	a.Post("/signin", func(ctx *fiber.Ctx) error {
		return controllers.SignInCustomer(ctx, customerClient)
	})
}
