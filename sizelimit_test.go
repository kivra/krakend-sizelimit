package sizelimit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIntNoUnit(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10"
	parsed := ParseMaxSize(cfg)
	if parsed != 10 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestFloatNoUnit(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10.0"
	parsed := ParseMaxSize(cfg)
	if parsed != 10 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestB(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10B"
	parsed := ParseMaxSize(cfg)
	if parsed != 10 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestBFloat(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10.7B"
	parsed := ParseMaxSize(cfg)
	if parsed != 10 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestKB(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10kB"
	parsed := ParseMaxSize(cfg)
	if parsed != 10000 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestKBFloat(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10.7kB"
	parsed := ParseMaxSize(cfg)
	if parsed != 10700 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestMB(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10MB"
	parsed := ParseMaxSize(cfg)
	if parsed != 10000000 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestGB(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10GB"
	parsed := ParseMaxSize(cfg)
	if parsed != 10000000000 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestTB(t *testing.T) {
	cfg := new(Config)
	cfg.MaxSize = "10TB"
	parsed := ParseMaxSize(cfg)
	if parsed != 10000000000000 {
		t.Fatalf("MaxSize parsed incorrectly. Got %d", parsed)
	}
}

func TestRequestSizeBelowLimit(t *testing.T) {
	body := []byte("hello")
	router := makeRouter(int64(len(body)) + 1)
	code := performRequest(body, router)

	if code != http.StatusOK {
		t.Fatalf("returned status %v. should return 200", code)
	}
}

func TestRequestSizeAtLimit(t *testing.T) {
	body := []byte("hello")
	router := makeRouter(int64(len(body)))
	code := performRequest(body, router)

	if code != http.StatusOK {
		t.Fatalf("returned status %v. should return 200", code)
	}
}

func TestRequestSizeAboveLimit(t *testing.T) {
	body := []byte("hello")
	router := makeRouter(int64(len(body)) - 1)
	code := performRequest(body, router)

	if code != http.StatusRequestEntityTooLarge {
		t.Fatalf("returned status %v. should return 413", code)
	}
}

func makeRouter(limit int64) *gin.Engine {
	router := gin.New()
	limiter := func(c *gin.Context) {
		LimiterFactory(limit, func(c *gin.Context) { c.Next() })(c)
	}
	success := func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	}
	router.POST("/test", limiter, success)
	return router
}

func performRequest(body []byte, router *gin.Engine) int {
	buf := new(bytes.Buffer)
	buf.Write(body)
	r := httptest.NewRequest(http.MethodPost, "/test", buf)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w.Code
}
