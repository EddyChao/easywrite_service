package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SwaggerMiddleware(swaggerUsername string, swaggerPassword string) func(c *gin.Context) {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if hasAuth && username == swaggerUsername && password == swaggerPassword {
			if c.Request.URL.Path == "/swagger/" {
				c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
				return
			}
			c.Next()
		} else {
			c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
