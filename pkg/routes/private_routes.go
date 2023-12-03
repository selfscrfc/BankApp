package routes

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func PrivateRoutes(a *fiber.App, customerService *grpc.ClientConn) {
	//a.Get("/account", middleware.JWTProtected())
	//a.Get("/t_history", middleware.JWTProtected())
	//a.Get("/user", middleware.JWTProtected(), func(ctx *fiber.Ctx) error { return controllers.GetUserDetails(ctx, customerService) })
	//
	//a.Post("/account", middleware.JWTProtected())
	//a.Post("/transaction", middleware.JWTProtected())
	//a.Post("/sign/out", middleware.JWTProtected())
	//
	//a.Delete("/account", middleware.JWTProtected())
	//a.Delete("/user", func(ctx *fiber.Ctx) error { return controllers.BlockUser(ctx, customerService) })
}
