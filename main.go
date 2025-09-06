package main

import (
	"NoLetServer/config"
	"NoLetServer/controller"
	"NoLetServer/database"
	"NoLetServer/model"
	"NoLetServer/push"
	"NoLetServer/router"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v3"
)

var (
	version   string
	buildDate string
	commitID  string
)

func main() {

	app := &cli.Command{
		Name:    "NowordsServer",
		Usage:   "Push Server For Nowords",
		Flags:   config.Flags(),
		Authors: []any{"to@uuneo.com"},
		Action: func(ctx context.Context, command *cli.Command) error {

			configPath := command.String("config")

			if _, err := os.Stat(configPath); err == nil {
				var conf model.Config
				ko := koanf.New(".")
				// Load JSON config.
				if err := ko.Load(file.Provider(configPath), yaml.Parser()); err != nil {
					log.Fatalf("error loading config: %v", err)
				} else {
					if err := ko.Unmarshal("", &conf); err != nil {
						log.Fatal(err)
					} else {
						config.LocalConfig.SetConfig(conf)
					}
				}

			}

			config.LocalConfig.System.Version = version
			config.LocalConfig.System.BuildDate = buildDate
			config.LocalConfig.System.CommitID = commitID

			// 创建 HTML 引擎
			systemConfig := config.LocalConfig.System

			engine := html.New("./static/html", ".html")

			fiberApp := fiber.New(fiber.Config{
				ServerHeader:          "NoLetServer",
				Concurrency:           systemConfig.Concurrency,
				ReadTimeout:           systemConfig.ReadTimeout,
				WriteTimeout:          systemConfig.WriteTimeout,
				IdleTimeout:           systemConfig.IdleTimeout,
				ProxyHeader:           systemConfig.ProxyHeader,
				ReduceMemoryUsage:     systemConfig.ReduceMemoryUsage,
				JSONEncoder:           sonic.Marshal,
				JSONDecoder:           sonic.Unmarshal,
				DisableStartupMessage: !systemConfig.Debug,
				Network:               "tcp",
				Views:                 engine,
				ErrorHandler: func(c *fiber.Ctx, err error) error {
					code := fiber.StatusInternalServerError
					var e *fiber.Error
					if errors.As(err, &e) {
						code = e.Code
					}
					return c.Status(code).JSON(model.BaseRes(code, err.Error()))
				},
			})

			fiberApp.Static(config.LocalConfig.System.URLPrefix, "static")

			// 监听结束信号
			MonitoringSignal(fiberApp)

			// 初始化数据库
			database.InitDatabase()

			fiberRouter := fiberApp.Group(config.LocalConfig.System.URLPrefix)

			router.SetupMiddler(fiberRouter, systemConfig.TimeZone)
			router.RegisterRoutes(fiberRouter)
			push.CreateAPNSClient(systemConfig.MaxAPNSClientCount)

			// 循环推送
			controller.CirclePush()
			controller.CirclePushPTT()
			controller.CircleDeleteExFile(systemConfig.Expired)

			if systemConfig.Cert != "" && systemConfig.Key != "" {
				return fiberApp.ListenTLS(systemConfig.Addr, systemConfig.Cert, systemConfig.Key)
			}
			return fiberApp.Listen(systemConfig.Addr)
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}

func MonitoringSignal(fiberApp *fiber.App) {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs // 等待一次
		log.Warnf("Received signal %s, shutting down...", sig)

		if err := fiberApp.Shutdown(); err != nil {
			log.Errorf("Server forced to shutdown error: %v", err)
		}
		if err := database.DB.Close(); err != nil {
			log.Errorf("Database close error: %v", err)
		}

		os.Exit(0)
	}()
}
