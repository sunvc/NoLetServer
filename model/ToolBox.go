package model

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sunvc/apns2"
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

func Unique[T comparable](list []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(list))

	for _, v := range list {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func InList[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
