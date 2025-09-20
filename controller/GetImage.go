package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
)

func GetImage(c *gin.Context) {
	fileName := c.Param("deviceKey")
	color := c.Query("color")

	path := filepath.Join("static", fileName)

	if fileName == "logo.svg" {
		c.Data(http.StatusOK, common.MIMEImageSvg, []byte(common.LogoSvgImage(color, true)))
		return
	}
	data, err := common.StaticFS.ReadFile(path)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, common.MIMEImagePng, data)

}
