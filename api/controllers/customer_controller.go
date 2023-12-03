package controllers

import (
	"PetBank/api/models"
	"PetBank/proto/Customers"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"time"
)

func CreateUser(c *fiber.Ctx, customerService *grpc.ClientConn) error {
	client := Customers.NewCustomerClient(customerService)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	customer := &models.Customer{}
	if err := c.BodyParser(customer); err != nil {
		return err
	}

	resp, err := client.Create(ctx, &Customers.CreateRequest{
		FullName: customer.FullName,
		Login:    customer.Login,
		Password: customer.Password,
	})

	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Error)
	}
	return nil
}

func GetUserDetails(c *fiber.Ctx, customerService *grpc.ClientConn) error {
	return nil
}

func BlockUser(c *fiber.Ctx, customerService *grpc.ClientConn) error {
	return nil
}
