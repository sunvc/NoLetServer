package router

import (
	"NoLetServer/controller"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {

	router.Get("/", controller.HomeController)
	router.Get("/info", controller.GetInfo)

	// 注册
	router.Post("/register", controller.RegisterController)
	router.Get("/register/:deviceKey", controller.RegisterController)

	router.Get("/ping", controller.Ping)
	router.Get("/p", controller.Ping)

	router.Get("/health", func(c *fiber.Ctx) error { return c.JSON("ok") })
	router.Get("/h", func(c *fiber.Ctx) error { return c.JSON("ok") })
	// 推送请求
	router.Post("/push", controller.BaseController)

	router.Get("/upload*", controller.UploadController)
	router.Get("/u", controller.UploadController)
	router.Post("/upload", controller.UploadController)
	router.Post("/u", controller.UploadController)
	router.Get("/image/:filename", controller.GetImage)
	router.Get("/img/:filename", controller.GetImage)

	group := router.Group("/ptt", Verification())
	{
		group.Post("/join", controller.JoinChannel)
		group.Post("/leave", controller.LeaveChannel)
		group.Get("/ping/:channel", controller.PingPTT)
		group.Post("/send", controller.UploadVoice)
		group.Get("/voices/:fileName", controller.GetVoice)
	}

	// title subtitle body
	router.Get("/:devicekey/:title/:subtitle/:body", controller.BaseController)
	router.Post("/:devicekey/:title/:subtitle/:body", controller.BaseController)
	// title body
	router.Get("/:devicekey/:title/:body", controller.BaseController)
	router.Post("/:deviceKey/:title/:body", controller.BaseController)
	// body
	router.Get("/:devicekey/:body", controller.BaseController)
	router.Post("/:devicekey/:body", controller.BaseController)
	// 参数化的推送
	router.Get("/:devicekey", controller.BaseController)
	router.Post("/:devicekey", controller.BaseController)

}
