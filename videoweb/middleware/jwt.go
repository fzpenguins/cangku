package middleware

import (
	"context"
	"net/http"
	"videoweb/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var code int
		code = http.StatusOK
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")
		if string(accessToken) == "" {
			code = http.StatusBadRequest
			c.JSON(http.StatusOK, map[string]interface{}{
				"status": code,
				"data":   "Token不能为空",
			})
			c.Abort()
			return
		}
		newAccessToken, newRefreshToken, err := utils.ParseRefreshToken(string(accessToken), string(refreshToken))
		if err != nil {
			code = http.StatusBadRequest
		}
		if code != http.StatusOK {
			c.JSON(http.StatusOK, map[string]interface{}{
				"status": code,
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		c.Header("access_token", newAccessToken)
		c.Header("refresh_token", newRefreshToken)
		c.Next(ctx)
	}
}
