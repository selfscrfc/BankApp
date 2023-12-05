package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
)

// CreateAccount	godoc
// @Summary		create
// @Tags Accounts
// @Accept json
// @Produce json
// @Param data body string true "Userid, IsCredit, Balance, Currency"
// @Success 200 {object} models.Account
// @Router /newaccount [post]
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

// GetAccountDetails	godoc
// @Summary		create
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param userid path string true "userid"
// @Success 200 {object} models.Account
// @Router /getaccountdetails [get]
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

// GetAllAccounts	godoc
// @Summary		create
// @Tags Accounts
// @Produce json
// @Success 200 {array} models.Account
// @Router /getallaccounts [get]
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

// BlockAccount	godoc
// @Summary		create
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param userid path string true "userid"
// @Success 200 {object} models.Account
// @Router /blockaccount [get]
func BlockAccount(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	acc := &Accounts.Account{}

	if err := c.BodyParser(acc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
		})
	}

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

// RW	godoc
// @Summary		create
// @Tags Accounts
// @Accept json
// @Produce json
// @Param data body string true "data"
// @Success 200 {object} models.Account
// @Router /balance [post]
func RW(c *fiber.Ctx, client *Accounts.AccountServiceClient, type_ bool) error {
	ctx := context.Background()

	req := &Accounts.RWRequest{}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
		})
	}

	resp, err := (*client).RW(ctx, &Accounts.RWRequest{
		AId:    req.AId,
		UId:    req.UId,
		Amount: req.Amount,
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     "successful",
		"account": resp,
	})
}
