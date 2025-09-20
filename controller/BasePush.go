package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
	"github.com/sunvc/NoLet/database"
	"github.com/sunvc/NoLet/push"
	"github.com/sunvc/apns2"
)

// BasePush 处理基础推送请求
// 验证推送参数并执行推送操作
func BasePush(c *gin.Context) {

	result := common.NewParamsResult(c)

	if len(result.Tokens) <= 0 {
		for _, key := range result.Keys {
			if len(key) > 5 {
				if token, err := database.DB.DeviceTokenByKey(key); err == nil {
					result.Tokens = append(result.Tokens, token)
				}

			}
		}
	}

	if len(result.Tokens) <= 0 {
		c.JSON(http.StatusOK, common.Failed(http.StatusBadRequest, "Failed to get device token"))
		return
	}

	if err := push.BatchPush(result, apns2.PushTypeAlert); err != nil {
		c.JSON(http.StatusOK, common.Failed(http.StatusInternalServerError, "push failed: %v", err))
		return
	}

	// 如果是管理员，加入到未推送列表
	if id, ok := result.Get("id").(string); common.Admin(c) && ok && len(id) > 0 {
		UpdateNotPushedData(id, result, apns2.PushTypeAlert)
	}

	c.JSON(http.StatusOK, common.Success())
}
