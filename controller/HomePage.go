package controller

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// HomeController 处理首页请求
// 支持两种功能:
// 1. 通过id参数移除未推送数据
// 2. 生成二维码图片
func HomeController(c *fiber.Ctx) error {

	if id := c.Query("id"); id != "" {
		RemoveNotPushedData(id)

		return c.SendStatus(http.StatusOK)
	}
	url := func() string {
		if code := c.Query("code"); code != "" {
			return code
		} else {

			if c.Secure() {
				return "https://" + string(c.Request().Host())
			}
			return "http://" + string(c.Request().Host())
		}
	}()
	params := c.Queries()

	params["URL"] = url
	params["ICP"] = config.LocalConfig.System.ICPInfo

	log.Info("logo:", params)

	return c.Render("index", params)
}

// Ping 处理心跳检测请求
// 返回服务器当前状态
func Ping(c *fiber.Ctx) error {
	return c.JSON(model.BaseRes(http.StatusOK, "pong"))
}
