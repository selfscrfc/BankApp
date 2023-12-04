package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	"github.com/selfscrfc/PetBank/utils"
	Customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
	"time"
)

// CreateUser	godoc
// @Summary		create
// @Tags Users
// @Accept json
// @Produce json
// @Param data body string true "data"
// @Param data body string true "data"
// @Param data body string true "data"
// @Success 200
// @Router /sign/up/ [post]
func CreateUser(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	sign := &models.SignUp{}

	if err := c.BodyParser(sign); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validator := utils.NewValidator()

	if err := validator.Struct(sign); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	customer := &models.Customer{
		Id:          uuid.New(),
		TimeCreated: time.Now(),
		FullName:    sign.FullName,
		Login:       sign.Login,
		Password:    utils.GeneratePassword(sign.Password),
		IsBlocked:   false,
	}

	_, err := (*client).Create(ctx, &Customers.CreateRequest{
		Id:       customer.Id.String(),
		Time:     customer.TimeCreated.Unix(),
		FullName: customer.FullName,
		Login:    customer.Login,
		Password: customer.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// GetUserDetails	godoc
// @Summary		Get user details with id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Router /user [get]
func GetUserDetails(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	customer := &models.Customer{}

	if err := c.BodyParser(customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	resp, err := (*client).GetDetails(ctx, &Customers.GetDetailsRequest{
		Id: customer.Id.String(),
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	customer.TimeCreated = time.Unix(resp.Time, 0)
	customer.FullName = resp.FullName
	customer.IsBlocked = resp.IsBlocked
	customer.Login = resp.Login

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// BlockUser	godoc
// @Summary		Block user by himself with id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Router /user/{id}/block [post]
func BlockUser(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	var id string

	if err := c.BodyParser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	resp, err := (*client).Block(ctx, &Customers.BlockRequest{
		BlockId: id,
	})

	if resp.Success {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}

// SignInUser	godoc
// @Summary		Sign in
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Router /sign/in [post]
func SignInUser(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	sign := &models.Sign{}
	validator := utils.NewValidator()

	customer := &models.Customer{}

	if err := c.BodyParser(customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	resp, err := (*client).Block(ctx, &Customers.BlockRequest{
		BlockId: customer.Id,
	})

	if resp.Success {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}
