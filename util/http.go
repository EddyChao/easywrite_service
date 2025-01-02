package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CopyResponseHeader(c *gin.Context, resp *http.Response) {
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
		c.Writer.Header().Add("Set-Cookie", cookie.String())
	}
	c.Writer.WriteHeader(resp.StatusCode)
	for k, vList := range resp.Header {
		if k != "Set-Cookie" {
			for _, v := range vList {
				c.Writer.Header().Add(k, v)
			}
		}
	}
}

func CopyRequestHeader(c *gin.Context, req *http.Request) {
	for k, vList := range c.Request.Header {
		newKey := k
		if strings.HasPrefix(k, "X-") {
			newKey = k[len("X-"):]
			c.Request.Header.Del(newKey)
			req.Header.Del(newKey)
		}
		for _, v := range vList {
			req.Header.Add(newKey, v)
		}
	}
}
