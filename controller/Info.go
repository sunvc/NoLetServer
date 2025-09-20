package controller

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
	"github.com/sunvc/NoLet/database"
)

func Info(c *gin.Context) {
	admin, ok := c.Get("admin")
	system := common.LocalConfig.System

	results := gin.H{
		"version": system.Version,
		"build":   system.BuildDate,
		"commit":  system.CommitID,
	}

	if ok && admin.(bool) {
		devices, _ := database.DB.CountAll()
		results["devices"] = devices
		results["arch"] = runtime.GOOS + "/" + runtime.GOARCH
		results["cpu"] = runtime.NumCPU()
	}
	c.JSON(http.StatusOK, results)
}
