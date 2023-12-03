package routes

import (
	"PetBank/api/controllers"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func PublicRoutes(a *fiber.App, customerService *grpc.ClientConn) {
	a.Post("/sign/up", func(ctx *fiber.Ctx) error { return controllers.CreateUser(ctx, customerService) })
	a.Post("/sign/in")
}
