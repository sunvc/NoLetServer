package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
)

// Ping 处理心跳检测请求
// 返回服务器当前状态
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, common.BaseResp{
		Code:      http.StatusOK,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	})
}

func Health(c *gin.Context) { c.String(http.StatusOK, "OK") }
