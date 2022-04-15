package sizelimit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/proxy"
	krakendgin "github.com/luraproject/lura/v2/router/gin"
)

func ExceedsSizeLimit(c *gin.Context, limit int64) bool {
	contentLength := c.Request.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(contentLength, 10, 64)
	if size > limit { // trust Content-Length header only if it exceeds MaxSize
		return true
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
	bodyBuffer := new(bytes.Buffer)
	_, err := io.Copy(bodyBuffer, c.Request.Body)
	c.Request.Body = io.NopCloser(bodyBuffer)
	return err != nil
}

func LimiterFactory(limit int64, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	respBody := map[string]interface{}{
		"code":          41300,
		"short_message": "Content Too Large",
		"long_message":  "Content should not exceed " + fmt.Sprint(limit) + " B",
	}

	return func(c *gin.Context) {
		if ExceedsSizeLimit(c, limit) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, respBody)
			return
		}
		handlerFunc(c)
	}
}

func HandlerFactory(next krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	return func(remote *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handlerFunc := next(remote, p)

		cfg, ok := ConfigGetter(remote.ExtraConfig)
		if !ok {
			return handlerFunc
		}

		return LimiterFactory(cfg.MaxSize, handlerFunc)
	}
}
