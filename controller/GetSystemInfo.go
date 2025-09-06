package controller

import (
	"NoLetServer/config"
	"NoLetServer/database"
	"NoLetServer/model"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func GetInfo(c *fiber.Ctx) error {
	system := config.LocalConfig.System
	devices, _ := database.DB.CountAll()

	if model.Admin(c) {
		return c.JSON(map[string]interface{}{
			"version": system.Version,
			"build":   system.BuildDate,
			"commit":  system.CommitID,
			"devices": devices,
			"arch":    runtime.GOOS + "/" + runtime.GOARCH,
			"cpu":     runtime.NumCPU(),
		})
	} else {
		return c.JSON(map[string]interface{}{
			"version": system.Version,
		})
	}

}
