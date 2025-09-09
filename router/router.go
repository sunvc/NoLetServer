package router

import (
	"NoLetServer/config"
	"NoLetServer/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func RegisterRoutes(router fiber.Router) {

	router.Get("/", controller.HomeController)
	router.Get("/info", controller.GetInfo)
	router.Get("/metrics", monitor.New(monitor.Config{Title: config.LocalConfig.System.Name}))

	// 注册
	router.Post("/register", CheckUserAgent, controller.RegisterController)
	router.Get("/register/:deviceKey<minLen(5);alpha>", CheckUserAgent, controller.RegisterController)

	router.Get("/ping", controller.Ping)

	router.Get("/health", func(c *fiber.Ctx) error { return c.JSON("ok") })
	// 推送请求
	router.Post("/push", controller.BaseController)

	router.Get("/upload", controller.UploadController)
	router.Post("/upload", controller.UploadController)
	router.Get("/image/:filename<minLen(3)>", controller.GetImage)
	router.Get("/img/:filename<minLen(3)>", controller.GetImage)

	{
		groupPtt := router.Group("/ptt", CheckUserAgent)
		groupPtt.Post("/join", controller.JoinChannel)
		groupPtt.Post("/leave", controller.LeaveChannel)
		groupPtt.Get("/ping/:channel<minLen(3)>", controller.PingPTT)
		groupPtt.Post("/send", controller.UploadVoice)
		groupPtt.Get("/voices/:fileName<minLen(3)>", controller.GetVoice)
	}

	// title subtitle body
	router.Get("/:devicekey<minLen(5);alpha>/:title/:subtitle/:body", controller.BaseController)
	router.Post("/:devicekey<minLen(5);alpha>/:title/:subtitle/:body", controller.BaseController)
	// title body
	router.Get("/:devicekey<minLen(5);alpha>/:title/:body", controller.BaseController)
	router.Post("/:deviceKey<minLen(5);alpha>/:title/:body", controller.BaseController)
	// body
	router.Get("/:devicekey<minLen(5);alpha>/:body", controller.BaseController)
	router.Post("/:devicekey<minLen(5);alpha>/:body", controller.BaseController)
	// 参数化的推送
	router.Get("/:devicekey<minLen(5);alpha>", controller.BaseController)
	router.Post("/:devicekey<minLen(5);alpha>", controller.BaseController)

}
