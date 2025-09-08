package controller

import (
	"NoLetServer/config"
	"NoLetServer/database"
	"NoLetServer/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// RegisterController 处理设备注册请求
// 支持 GET 和 POST 两种请求方式:
// GET: 检查设备key是否存在
// POST: 注册新的设备token
func RegisterController(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		deviceKey := c.Params("deviceKey")
		if deviceKey == "" {
			return c.JSON(model.Failed(http.StatusBadRequest, "device key is empty"))
		}
		if database.DB.KeyExists(deviceKey) {
			return c.JSON(model.Success())
		} else {
			auth := model.Admin(c)

			if auth {
				_, err := database.DB.SaveDeviceTokenByKey(deviceKey, "")
				if err != nil {

					return c.JSON(model.Failed(http.StatusBadRequest, "device key is not exist"))
				}

				return c.JSON(model.Success())
			}

			return c.JSON(model.Failed(http.StatusBadRequest, "device key is not exist"))
		}
	}

	var err error
	var device = new(model.DeviceInfo)
	device.Voice = config.LocalConfig.System.Voice

	if err = c.BodyParser(&device); err != nil {
		return c.JSON(model.Failed(http.StatusBadRequest, "failed to get device token: %v", err))
	}

	if device.Token == "" {
		return c.JSON(model.Failed(http.StatusBadRequest, "deviceToken is empty"))
	}

	device.Key, err = database.DB.SaveDeviceTokenByKey(device.Key, device.Token)

	if err != nil {
		return c.JSON(model.Failed(http.StatusInternalServerError, "device registration failed: %v", err))
	}

	return c.JSON(model.Success(device))
}
