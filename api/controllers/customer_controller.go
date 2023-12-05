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

// CreateCustomer	godoc
// @Summary		create
// @Tags Users
// @Accept json
// @Produce json
// @Param data body string true "data"
// @Success 200 {object} models.Customer
// @Router /signup [post]
func CreateCustomer(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	sign := &models.SignUp{}

	if err := c.BodyParser(sign); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
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
		FullName: sign.FullName,
		Login:    sign.Login,
		Password: sign.Password,
	}

	resp, err := (*client).Create(ctx, &Customers.CreateRequest{
		FullName: customer.FullName,
		Login:    customer.Login,
		Password: customer.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	customer.Id, err = uuid.Parse(resp.Id)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "uuid parse error: " + err.Error(),
		})
	}
	customer.TimeCreated = time.Unix(resp.Time, 0)
	customer.IsBlocked = false
	customer.Password = ""

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// GetCustomerDetails	godoc
// @Summary		Get user details with id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Customer
// @Router /user/{id} [get]
func GetCustomerDetails(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	id := c.Params("id")

	resp, err := (*client).GetDetails(ctx, &Customers.GetDetailsRequest{
		Id: id,
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	id_, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	customer := &models.Customer{Id: id_}
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

// BlockCustomer	godoc
// @Summary		Block user by himself with id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200
// @Router /user/{id}/block [get]
func BlockCustomer(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	id := c.Params("id")

	_, err := (*client).Block(ctx, &Customers.BlockRequest{
		BlockId: id,
	})

	if err != nil {
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

// SignInCustomer	godoc
// @Summary		Sign in
// @Tags Users
// @Accept json
// @Produce json
// @Param id body string true "id"
// @Success 200 {object} models.Customer
// @Router /signin [post]
func SignInCustomer(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	customer := &models.Customer{}

	if err := c.BodyParser(customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Parse error: " + err.Error(),
		})
	}

	resp, err := (*client).SignIn(ctx, &Customers.SignInRequest{
		Login:    customer.Login,
		Password: customer.Password,
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Authorization error: " + err.Error(),
		})
	}

	customer.TimeCreated = time.Unix(resp.Time, 0)
	customer.Id, err = uuid.Parse(resp.Id)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	customer.FullName = resp.FullName
	customer.Login = resp.Login
	customer.IsBlocked = resp.IsBlocked

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": customer,
	})
}

// GetAllCustomers	godoc
// @Summary	get all
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.Customer
// @Router /getallusers [get]
func GetAllCustomers(c *fiber.Ctx, client *Customers.CustomerClient) error {
	ctx := context.Background()

	var resp *Customers.GetAllResponse

	resp, err := (*client).GetAll(ctx, &Customers.GetAllRequest{})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Authorizatione error: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":     false,
		"msg":       nil,
		"customers": resp,
	})
}
