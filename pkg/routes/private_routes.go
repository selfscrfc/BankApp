package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/selfscrfc/PetBank/api/controllers"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

func PrivateRoutes(a *fiber.App, customerClient *customers.CustomerClient,
	accountsClient *account.AccountServiceClient) {

	a.Get("/account", func(ctx *fiber.Ctx) error {
		return controllers.GetAccountDetails(ctx, accountsClient)
	})
	a.Get("/getallaccounts", func(ctx *fiber.Ctx) error {
		return controllers.GetAllAccounts(ctx, accountsClient)
	})
	a.Get("/user", func(ctx *fiber.Ctx) error {
		return controllers.GetCustomerDetails(ctx, customerClient)
	})
	a.Get("/getallusers", func(ctx *fiber.Ctx) error {
		return controllers.GetAllCustomers(ctx, customerClient)
	})

	a.Post("/newaccount", func(ctx *fiber.Ctx) error {
		return controllers.CreateAccount(ctx, accountsClient)
	})

	a.Delete("/account", func(ctx *fiber.Ctx) error {
		return controllers.BlockAccount(ctx, accountsClient)
	})
	a.Delete("/user/{id}/block", func(ctx *fiber.Ctx) error {
		return controllers.BlockCustomer(ctx, customerClient)
	})
}
