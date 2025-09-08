package router

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"encoding/base64"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	authFreeRouters = []string{"/ping", "/register", "/health", "/u", "/upload", "/image", "/img", "/ptt"}
)

func SetupMiddler(router fiber.Router, timeZone string) {
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
	router.Use("/register/+", CheckUserAgent())
	router.Use(favicon.New(favicon.Config{
		File: "./static/logo.svg",
	}))
}

func AuthRouter() fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		ctx.Locals("admin", false)

		urlPrefix := config.LocalConfig.System.URLPrefix
		// Get authorization header
		auth := ctx.Get(fiber.HeaderAuthorization)

		for _, item := range authFreeRouters {
			if strings.HasPrefix(ctx.Path(), path.Join(urlPrefix, item)) {
				return ctx.Next()
			}
		}

		user := config.LocalConfig.System.User
		password := config.LocalConfig.System.Password

		if user == "" || password == "" {
			return ctx.Next()
		}

		// Check if the header contains content besides "basic".
		if len(auth) < 6 || !strings.HasPrefix(strings.ToLower(auth), "basic ") {
			return ctx.Status(401).SendString("I'm a teapot")
		}

		// Decode the header contents
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return ctx.Status(401).SendString("I'm a teapot")
		}

		// Get the credentials
		credentials := utils.UnsafeString(raw)

		// Check if the credentials are in the correct form
		// which is "username:password".
		index := strings.Index(credentials, ":")
		if index == -1 {
			return ctx.Status(401).SendString("I'm a teapot")
		}

		// Get the username and password
		username := credentials[:index]
		password2 := credentials[index+1:]
		if user == username && password == password2 {
			ctx.Locals("admin", true)
			return ctx.Next()
		}

		return ctx.Status(401).SendString("I'm a teapot")
	}
}

func CheckUserAgent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Info(ctx.Path())
		userAgent := ctx.Get(fiber.HeaderUserAgent)
		if strings.HasPrefix(userAgent, config.LocalConfig.System.Name) {
			return ctx.Status(401).SendString("I'm a teapot")
		}
		return ctx.Next()
	}
}

func Verification() fiber.Handler {

	return func(c *fiber.Ctx) error {
		log.Info(c.Get("X-Signature"))
		// 先查看是否是管理员身份
		authHeader := c.Get(fiber.HeaderAuthorization)
		UserAgent := c.Get(fiber.HeaderUserAgent)
		if !config.LocalConfig.System.Debug {
			if len(authHeader) < 10 || !strings.HasPrefix(UserAgent, config.LocalConfig.System.Name) {
				return c.Status(401).JSON(model.Failed(-1, "Unauthorized access"))
			}
		}

		return c.Next()

	}
}
