package router

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunvc/NoLet/common"
	"github.com/sunvc/NoLet/controller"
)

func Verification() gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodPost {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		// 先查看是否是管理员身份
		authHeader := c.GetHeader("Authorization")
		if common.Contains[string](common.LocalConfig.System.Auths, authHeader) && authHeader != "" {
			c.Set("admin", true)
			return
		}

		localUser := common.LocalConfig.System.User
		localPassword := common.LocalConfig.System.Password
		// 配置了账号密码，进行身份校验
		if localUser != "" && localPassword != "" {
			// 优先使用 Basic Auth
			user, pass, hasAuth := c.Request.BasicAuth()
			if !hasAuth {
				// 如果没有 Basic Auth，则尝试从查询参数中获取
				user = c.Query(common.UserName)
				pass = c.Query(common.Password)

				if c.Request.Method == http.MethodPost {
					if user == "" {
						user = c.PostForm(common.UserName)
					}
					if pass == "" {
						pass = c.PostForm(common.Password)
					}
				}
			}

			if user == localUser && pass == localPassword {
				c.Set("admin", true)
				return
			}

		}

		// 如果没有身份验证信息
		c.Set("admin", false)
		c.Next()
	}
}

// CheckDotParamMiddleware 检查 GET 请求第一个 path 参数是否包含 '.'
func CheckDotParamMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value := c.Param("deviceKey"); strings.Contains(value, ".") {
			controller.GetImage(c)
			c.Abort()
			return
		}
		// 放行请求
		c.Next()
	}
}

func CheckUserAgent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println(ctx.Request.URL.String())
		userAgent := ctx.GetHeader(common.HeaderUserAgent)
		if !strings.HasPrefix(userAgent, common.LocalConfig.System.Name) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
