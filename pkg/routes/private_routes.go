package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/selfscrfc/PetBank/api/controllers"
	"github.com/selfscrfc/PetBank/pkg/middleware"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

func PrivateRoutes(a *fiber.App, customerClient *customers.CustomerClient,
	accountsClient *account.AccountServiceClient) {

	//a.Get("/account", middleware.JWTProtected())
	//a.Get("/t_history", middleware.JWTProtected())
	a.Get("/user", middleware.JWTProtected(), func(ctx *fiber.Ctx) error {
		return controllers.GetUserDetails(ctx, customerClient)
	})
	//
	//a.Post("/account", middleware.JWTProtected())
	//a.Post("/transaction", middleware.JWTProtected())
	//a.Post("/sign/out", middleware.JWTProtected())
	//
	//a.Delete("/account", middleware.JWTProtected())
	a.Delete("/user/{id}/block", middleware.JWTProtected(), func(ctx *fiber.Ctx) error {
		return controllers.BlockUser(ctx, customerClient)
	})
}
