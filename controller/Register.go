package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
	"github.com/sunvc/NoLet/database"
)

// Register 处理设备注册请求
// 支持 GET 和 POST 两种请求方式:
// GET: 检查设备key是否存在
// POST: 注册新的设备token
func Register(c *gin.Context) {
	if c.Request.Method == "GET" {
		deviceKey := c.Param("deviceKey")
		if deviceKey == "" {
			c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "device key is empty"))
			return
		}
		if database.DB.KeyExists(deviceKey) {
			c.JSON(http.StatusOK, common.Success())
			return
		} else {
			admin, ok := c.Get("admin")

			if ok && admin.(bool) {
				_, err := database.DB.SaveDeviceTokenByKey(deviceKey, "")
				if err != nil {
					c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "device key is not exist"))
					return
				}
				c.JSON(http.StatusOK, common.Success())
				return
			}

			c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "device key is not exist"))
		}
		return
	}

	var err error
	var device common.DeviceInfo

	if err = c.BindJSON(&device); err != nil {
		c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "failed to get device token: %v", err))
		return
	}

	if device.Token == "" {
		c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "deviceToken is empty"))
		return
	}

	device.Key, err = database.DB.SaveDeviceTokenByKey(device.Key, device.Token)

	if err != nil {
		c.JSON(http.StatusOK, common.Failed(http.StatusInternalServerError, "device registration failed: %v", err))
	}

	c.JSON(http.StatusOK, common.Success(device))
}
