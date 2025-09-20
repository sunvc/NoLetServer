package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
	"github.com/sunvc/NoLet/database"
	"github.com/sunvc/NoLet/push"
	"github.com/sunvc/NoLet/router"
	"github.com/urfave/cli/v3"
)

var (
	version   string
	buildDate string
	commitID  string
)

//go:embed static/*
var staticFS embed.FS

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctxOut, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	common.StaticFS = &staticFS
	app := &cli.Command{
		Name:    "NoLetServer",
		Usage:   "Push Server For NoLet",
		Flags:   common.Flags(),
		Authors: []any{"to@uuneo.com"},
		Action: func(_ context.Context, command *cli.Command) error {

			if configPath := command.String("config"); configPath != "" {
				common.LocalConfig.SetConfig(configPath)
			}

			common.LocalConfig.System.Version = version
			common.LocalConfig.System.BuildDate = buildDate
			common.LocalConfig.System.CommitID = commitID

			database.InitDatabase()

			systemConfig := common.LocalConfig.System

			if systemConfig.Debug {
				gin.SetMode(gin.DebugMode)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}
			gin.ForceConsoleColor()

			engine := gin.Default()
			engine.Use(router.Verification())

			tmpl := template.Must(template.New("").ParseFS(staticFS, "static/*.html"))
			engine.SetHTMLTemplate(tmpl)

			push.CreateAPNSClient(systemConfig.MaxAPNSClientCount)

			router.SetupRouter(engine)

			var tLSConfig *tls.Config

			if systemConfig.Key != "" && systemConfig.Cert != "" {
				cert, err := tls.LoadX509KeyPair(systemConfig.Key, systemConfig.Cert)
				if err == nil {
					tLSConfig = &tls.Config{
						Certificates: []tls.Certificate{cert},
						MinVersion:   tls.VersionTLS12,
					}
				}

			}

			server := &http.Server{
				Addr:           systemConfig.Addr,
				Handler:        engine,
				TLSConfig:      tLSConfig,
				ReadTimeout:    systemConfig.ReadTimeout,
				WriteTimeout:   systemConfig.WriteTimeout,
				IdleTimeout:    systemConfig.IdleTimeout,
				MaxHeaderBytes: 1 << 12,
			}
			httpServerError := make(chan error, 1)
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					httpServerError <- err
					return
				}
				ctxOut.Done()
			}()

			select {
			case <-ctxOut.Done():
				log.Println("Received shutdown signal")
			case e := <-httpServerError:
				log.Printf("Server start error: %v", e)
			}

			ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// 关闭HTTP服务器
			if err := server.Shutdown(ctxShutdown); err != nil {
				log.Printf("Server forced to shutdown error: %v", err)
			}

			// 关闭数据库连接
			if err := database.DB.Close(); err != nil {
				log.Printf("Database close error: %v", err)
			}

			// 关闭APNS客户端资源
			push.CloseAPNSClients()

			log.Println("All resources have been properly released")
			return nil
		},
	}

	_ = app.Run(context.Background(), os.Args)
}
