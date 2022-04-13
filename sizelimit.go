package sizelimit

import (
	"bytes"
	"encoding/json"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/proxy"
	krakendgin "github.com/luraproject/lura/v2/router/gin"
)

const Namespace = "kivra/sizelimit"

type Config struct {
	MaxBytes int64 `json:"max_bytes"`
}

func ConfigGetter(e config.ExtraConfig) (*Config, bool) {
	cfg := new(Config)

	tmp, ok := e[Namespace]
	if !ok {
		return cfg, false
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tmp); err != nil {
		panic("sizelimit: Error: failed to parse config")
	}
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		panic("sizelimit: Error: failed to parse config")
	}

	return cfg, true
}

func HandlerFactory(next krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	return func(remote *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handlerFunc := next(remote, p)

		cfg, ok := ConfigGetter(remote.ExtraConfig)
		if !ok {
			return handlerFunc
		}

		limiter := limits.RequestSizeLimiter(cfg.MaxBytes)

		return func(c *gin.Context) {
			limiter(c)
			handlerFunc(c)
		}
	}
}
