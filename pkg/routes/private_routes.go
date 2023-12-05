package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/selfscrfc/PetBank/api/controllers"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

func PrivateRoutes(a *fiber.App, customerClient *customers.CustomerClient,
	accountsClient *account.AccountServiceClient) {

	a.Get("/getaccountdetails/:userid/:id", func(ctx *fiber.Ctx) error {
		return controllers.GetAccountDetails(ctx, accountsClient)
	})
	a.Get("/getallaccounts", func(ctx *fiber.Ctx) error {
		return controllers.GetAllAccounts(ctx, accountsClient)
	})
	a.Post("/newaccount", func(ctx *fiber.Ctx) error {
		return controllers.CreateAccount(ctx, accountsClient)
	})
	a.Post("/blockaccount", func(ctx *fiber.Ctx) error {
		return controllers.BlockAccount(ctx, accountsClient)
	})
	a.Post("/balance", func(ctx *fiber.Ctx) error {
		return controllers.RW(ctx, accountsClient)
	})

	a.Get("/user/:id", func(ctx *fiber.Ctx) error {
		return controllers.GetCustomerDetails(ctx, customerClient)
	})
	a.Get("/getallusers", func(ctx *fiber.Ctx) error {
		return controllers.GetAllCustomers(ctx, customerClient)
	})
	a.Get("/user/block/:id", func(ctx *fiber.Ctx) error {
		return controllers.BlockCustomer(ctx, customerClient)
	})
}
