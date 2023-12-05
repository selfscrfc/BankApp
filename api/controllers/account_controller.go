package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
)

func CreateAccount(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	acc := &models.Account{}

	if err := c.BodyParser(acc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
		})
	}

	resp, err := (*client).Create(ctx, &Accounts.CreateRequest{
		Id:       acc.Id.String(),
		UserId:   acc.UserId.String(),
		IsCredit: acc.IsCredit,
		Balance:  acc.Balance,
		Currency: Accounts.Currency(acc.Currency),
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	acc.Id, err = uuid.Parse(resp.Id)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "uuid parse error: " + err.Error(),
		})
	}
	acc.UserId, err = uuid.Parse(resp.UserId)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "uuid parse error: " + err.Error(),
		})
	}

	acc.IsCredit = resp.IsCredit
	acc.Balance = resp.Balance
	acc.Currency = models.Currency(resp.Currency)

	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"account": acc,
	})
}

func GetAccountDetails(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	acc := &models.Account{}

	if err := c.BodyParser(acc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
		})
	}

	resp, err := (*client).GetDetail(ctx, &Accounts.GetDetailsRequest{
		Id:     acc.Id.String(),
		UserId: acc.UserId.String(),
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	acc.IsCredit = resp.IsCredit
	acc.Balance = resp.Balance
	acc.Currency = models.Currency(resp.Currency)

	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"account": acc,
	})
}

func GetAllAccounts(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	resp, err := (*client).GetAll(ctx, &Accounts.GetAllRequest{})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	accs := resp.Accounts

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"accounts": accs,
	})
}

func BlockAccount(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	acc := &Accounts.Account{}

	resp, err := (*client).Block(ctx, &Accounts.BlockRequest{
		Id:     acc.Id,
		UserId: acc.UserId,
	})

	if resp.Success {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}
