package controller

import (
	"NoLetServer/database"
	"NoLetServer/model"
	"NoLetServer/push"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sunvc/apns2"
)

func BaseController(c *fiber.Ctx) error {
	result := model.NewParamsResult(c)

	if len(result.Tokens) <= 0 {
		for _, key := range result.Keys {
			if len(key) > 5 {
				if token, err := database.DB.DeviceTokenByKey(key); err == nil {
					result.Tokens = append(result.Tokens, token)
				}

			}
		}
	}

	if len(result.Tokens) <= 0 {
		return c.JSON(model.Failed(http.StatusBadRequest, "Failed to get device token"))
	}

	if err := push.BatchPush(result, apns2.PushTypeAlert); err != nil {
		return c.JSON(model.Failed(http.StatusInternalServerError, "push failed: %v", err))
	}

	// 如果是管理员，加入到未推送列表
	if id, ok := result.Get("id").(string); model.Admin(c) && ok && len(id) > 0 {
		UpdateNotPushedData(id, result, apns2.PushTypeAlert)
	}

	return c.JSON(model.Success())
}
