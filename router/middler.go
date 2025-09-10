package router

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	authFreeRouters = []string{"/ping", "/register", "/health", "/u", "/upload", "/image", "/img", "/ptt"}
)

func SetupMiddler(router *fiber.App, timeZone string) {
	router.Use(logger.New(logger.Config{
		Format:     "${time} INFO  ${ip} -> [${status}] ${method} ${latency}   ${route} => ${url} ${error} ${UserAgent}\n",
		TimeFormat: "2006-01-02 15:04:05",
		CustomTags: map[string]logger.LogFunc{
			"UserAgent": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				auth := c.Get(fiber.HeaderUserAgent)
				return output.Write([]byte(auth))
			},
		},
		Output:   os.Stdout,
		TimeZone: timeZone,
	}))
	router.Use(recover.New())
	router.Use(AuthRouter())
	router.Use(favicon.New(favicon.Config{
		File: "./static/logo.svg",
	}))

}

func AuthRouter() fiber.Handler {

	return basicauth.New(basicauth.Config{
		Next: func(ctx *fiber.Ctx) bool {
			ctx.Locals("admin", false)
			auth := ctx.Get(fiber.HeaderAuthorization)
			if auth == "123" {
				ctx.Locals("admin", true)
			}
			sysConfig := config.LocalConfig.System

			if sysConfig.User == "" || sysConfig.Password == "" {
				return true
			}

			for _, item := range authFreeRouters {
				if strings.HasPrefix(ctx.Path(), path.Join(sysConfig.URLPrefix, item)) {
					return true
				}
			}
			return false
		},
		Users: map[string]string{
			config.LocalConfig.System.User: config.LocalConfig.System.Password,
		},
		Realm: "Forbidden",
		Authorizer: func(s string, s2 string) bool {
			return s == config.LocalConfig.System.User && s2 == config.LocalConfig.System.Password
		},
		Unauthorized: func(ctx *fiber.Ctx) error {
			return ctx.JSON(model.Failed(-1, "Unauthorized access"))
		},
	})
}

func CheckUserAgent(ctx *fiber.Ctx) error {
	log.Info(ctx.Path())
	userAgent := ctx.Get(fiber.HeaderUserAgent)
	if !strings.HasPrefix(userAgent, config.LocalConfig.System.Name) {
		return ctx.Status(401).SendString("I'm a teapot")
	}
	return ctx.Next()
}
