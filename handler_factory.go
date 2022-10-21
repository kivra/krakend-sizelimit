package sizelimit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apierrors "github.com/kivra/kivra-api-errors"
	correlationid "github.com/kivra/krakend-correlationid"
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
	apierrors.Load()
	apiError := apierrors.FromStatusOrFallback(http.StatusRequestEntityTooLarge)
	apiError.Payload.LongMessage = fmt.Sprintf("Content length should not exceed %d B", limit)

	return func(c *gin.Context) {
		if ExceedsSizeLimit(c, limit) {
			c.Writer.Header().Set(apierrors.ErrorCodeHeader, apiError.Payload.Code)
			c.Writer.Header().Set(correlationid.Header, c.Request.Header.Get(correlationid.Header))
			c.AbortWithStatusJSON(apiError.StatusCode, apiError.Payload)
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
