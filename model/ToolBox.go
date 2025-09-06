package model

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uuneo/apns2"
)

func Admin(ctx *fiber.Ctx) bool {
	admin, ok := ctx.Locals("admin").(bool)
	return ok && admin
}

type NotPushedData struct {
	ID           string          `json:"id"`
	CreateDate   time.Time       `json:"createDate"`
	LastPushDate time.Time       `json:"lastPushDate"`
	Count        int             `json:"count"`
	Params       *ParamsResult   `json:"params"`
	PushType     apns2.EPushType `json:"pushType"`
}
