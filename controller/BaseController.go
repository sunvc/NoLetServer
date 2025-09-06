package controller

import (
	"NoLetServer/database"
	"NoLetServer/model"
	"NoLetServer/push"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/uuneo/apns2"
)

func BaseController(c *fiber.Ctx) error {
	result := model.NewParamsResult(c)

	if value, ok := result.Params.Get(model.DeviceToken); !ok || len(fmt.Sprint(value)) <= 0 {
		deviceValue, deviceOk := result.Params.Get(model.DeviceKey)
		if !deviceOk || len(fmt.Sprint(deviceValue)) <= 0 {
			return c.JSON(model.Failed(http.StatusBadRequest, "failed to get device token: deviceKey or deviceToken is required"))
		}
		token, err := database.DB.DeviceTokenByKey(fmt.Sprint(deviceValue))
		if err != nil {
			return c.JSON(model.Failed(http.StatusBadRequest, "failed to get device token: %v", err))
		}

		result.DeviceTokens = append(result.DeviceTokens, token)
	}
	if err := push.BatchPush(result, apns2.PushTypeAlert); err != nil {
		return c.JSON(model.Failed(http.StatusInternalServerError, "push failed: %v", err))
	}

	// 如果是管理员，加入到未推送列表
	if model.Admin(c) {
		UpdateNotPushedData(result.Get("id").(string), result, apns2.PushTypeAlert)
	}

	return c.JSON(model.Success())
}
