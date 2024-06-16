package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 允许跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*,Content-Type,Accept,Cache-Control,Origin,Referer,User-Agent,Content-Length,Host,Accept-Language,Connection,Accept-Encoding,Access-Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "*,Access-Token,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}
