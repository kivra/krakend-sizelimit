package sizelimit

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/luraproject/lura/v2/config"
)

const Namespace = "kivra/sizelimit"

type Config struct {
	MaxSize string `json:"max_size"`
}

type ConfigInternal struct {
	MaxSize int64 `json:"max_size"`
}

func ConfigGetter(e config.ExtraConfig) (*ConfigInternal, bool) {
	tmp, ok := e[Namespace]
	if !ok {
		return new(ConfigInternal), false
	}

	cfg := new(Config)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tmp); err != nil {
		panic("sizelimit: Error: failed to parse config")
	}
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		panic("sizelimit: Error: failed to parse config")
	}

	cfgInternal := new(ConfigInternal)
	cfgInternal.MaxSize = ParseMaxSize(cfg)

	return cfgInternal, true
}

func ParseMaxSize(cfg *Config) int64 {
	rex := regexp.MustCompile(`^([\d.]+)([a-zA-Z]*)$`)
	matches := rex.FindStringSubmatch(cfg.MaxSize)

	if matches == nil || len(matches) != 3 || matches[0] != cfg.MaxSize {
		panic("sizelimit: Error: invalid value for 'max_size'")
	}

	value, err := strconv.ParseFloat(matches[1], 64)

	if err != nil {
		panic("sizelimit: Error: invalid value for 'max_size'")
	}

	var multiplier float64
	switch matches[2] {
	case "":
		multiplier = 1
	case "B":
		multiplier = 1
	case "kB":
		multiplier = 1000
	case "MB":
		multiplier = 1000000
	case "GB":
		multiplier = 1000000000
	case "TB":
		multiplier = 1000000000000
	default:
		panic("sizelimit: Error: unknown unit: " + cfg.MaxSize)
	}
	return int64(value * multiplier)
}
