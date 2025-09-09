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
	router.Get("/register/:deviceKey", CheckUserAgent, controller.RegisterController)

	router.Get("/ping", controller.Ping)

	router.Get("/health", func(c *fiber.Ctx) error { return c.JSON("ok") })
	// 推送请求
	router.Post("/push", controller.BaseController)

	router.Get("/upload*", controller.UploadController)
	router.Post("/upload", controller.UploadController)
	router.Get("/image/:filename", controller.GetImage)
	router.Get("/img/:filename", controller.GetImage)

	// title subtitle body
	router.Get("/:devicekey<minLen(5)>/:title/:subtitle/:body", controller.BaseController)
	router.Post("/:devicekey<minLen(5)>/:title/:subtitle/:body", controller.BaseController)
	// title body
	router.Get("/:devicekey<minLen(5)>/:title/:body", controller.BaseController)
	router.Post("/:deviceKey<minLen(5)>/:title/:body", controller.BaseController)
	// body
	router.Get("/:devicekey<minLen(5)>/:body", controller.BaseController)
	router.Post("/:devicekey<minLen(5)>/:body", controller.BaseController)
	// 参数化的推送
	router.Get("/:devicekey<minLen(5)>", controller.BaseController)
	router.Post("/:devicekey<minLen(5)>", controller.BaseController)

}
