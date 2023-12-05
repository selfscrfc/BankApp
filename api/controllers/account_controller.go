package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
)

// CreateAccount	godoc
// @Summary	create new account
// @Tags Accounts
// @Accept json
// @Produce json
// @Param data body string true "Userid (string) , IsCredit (bool) , Balance (int) , Currency (1-3) | {"userid":"","iscredit":"","balance":"","currency":""}"
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
// @Summary get account details
// @Tags Accounts
// @Accept json
// @Produce json
// @Param userid path string true "userid uuid string"
// @Param id path string true "id uuid string"
// @Success 200 {object} models.Account
// @Router /getaccountdetails/{userid}/{id} [get]
func GetAccountDetails(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()
	id_, err := uuid.Parse(c.Params("id"))
	uid_, err := uuid.Parse(c.Params("userid"))

	acc := &models.Account{
		Id:     id_,
		UserId: uid_,
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
// @Summary	get all accounts
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
// @Summary		block acc from operations
// @Tags Accounts
// @Accept json
// @Produce json
// @Param data body string true "id (uuid string) , userid (uuid string) | {"id":"","userid":""}"
// @Success 200 {object} models.Account
// @Router /blockaccount [post]
func BlockAccount(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
	ctx := context.Background()

	acc := &Accounts.Account{}

	if err := c.BodyParser(acc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "parse error: " + err.Error(),
		})
	}
	_, err := (*client).Block(ctx, &Accounts.BlockRequest{
		Id:     acc.Id,
		UserId: acc.UserId,
	})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "rpc response error: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "block successful",
	})
}

// RW	godoc
// @Summary	replenish / withdraw account
// @Tags Accounts
// @Accept json
// @Produce json
// @Param data body string true "accountId (uuid string) , userid (uuid string) , amount (int) | {"aId":"","uId":"","amount":}"
// @Success 200 {object} models.Account
// @Router /balance [post]
func RW(c *fiber.Ctx, client *Accounts.AccountServiceClient) error {
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
		"msg":     "operation successful",
		"account": resp,
	})
}
