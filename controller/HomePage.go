package controller

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"html/template"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

	params := map[string]interface{}{
		"LOGOSVG": template.URL(config.LOGOSVG),
		"LOGOPNG": template.URL(config.LOGOPNG),
		"ICP":     config.LocalConfig.System.ICPInfo,
		"URL":     template.URL(url),
	}
	return c.Render("index", params)
}

// Ping 处理心跳检测请求
// 返回服务器当前状态
func Ping(c *fiber.Ctx) error {
	return c.JSON(model.BaseRes(http.StatusOK, "pong"))
}
