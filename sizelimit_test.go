package sizelimit

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apierrors "github.com/kivra/kivra-api-errors"
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
	router := makeRouter(int64(len(body))+1, body)
	res := performRequest(body, router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v '%s'. should return 200", res.Code, res.Body)
	}
}

func TestRequestSizeAtLimit(t *testing.T) {
	body := []byte("hello")
	router := makeRouter(int64(len(body)), body)
	res := performRequest(body, router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v %s. should return 200", res.Code, res.Body)
	}
}

func TestRequestSizeAboveLimit(t *testing.T) {
	body := []byte("hello")
	router := makeRouter(int64(len(body))-1, body)
	res := performRequest(body, router)

	if res.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("returned %v %s. should return 413", res.Code, res.Body)
	}

	errorCode := res.Header().Get(apierrors.ErrorCodeHeader)
	if errorCode != "41300" {
		t.Fatalf("returned %s. should return 41300", errorCode)
	}
}

func makeRouter(limit int64, bodySent []byte) *gin.Engine {
	router := gin.New()
	limiter := func(c *gin.Context) {
		LimiterFactory(limit, func(c *gin.Context) { c.Next() })(c)
	}
	success := func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		if !bytes.Equal(body, bodySent) {
			c.String(http.StatusInternalServerError, string(body))
		} else {
			c.String(http.StatusOK, "OK")
		}
	}
	router.POST("/test", limiter, success)
	return router
}

func performRequest(body []byte, router *gin.Engine) *httptest.ResponseRecorder {
	buf := new(bytes.Buffer)
	buf.Write(body)
	r := httptest.NewRequest(http.MethodPost, "/test", buf)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
