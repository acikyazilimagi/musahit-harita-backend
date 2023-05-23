package cache

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCache(t *testing.T) {
	// Create new Fiber instance and use the New middleware
	app := fiber.New()
	app.Use(pprof.New())
	app.Use(New())
	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/swagger/*"
		},
		Level: compress.LevelBestCompression,
	}))

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Response().SetStatusCode(200)
		return c.Send([]byte(`{"test": "test"}`))
	})
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		// set body
		return c.SendStatus(200)
	})
	app.Get("/metrics", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Case when path is /healthcheck
	t.Run("GET HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:2222/healthcheck", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Case when path is /metrics
	t.Run("GET Metrics", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:2222/metrics", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Case when path is /test
	t.Run("GET Test", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Case when path is /test
	t.Run("GET Test should be cached", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)

		// Check if response is cached
		cachedResponse := resp.Header.Get("x-cached-response")
		assert.Equal(t, "true", cachedResponse)
	})

	// Case when path is /test
	t.Run("POST Test", func(t *testing.T) {
		req := httptest.NewRequest("POST", "http://localhost:2222/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Case when path is /test
	t.Run("DELETE Test", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "http://localhost:2222/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 405, resp.StatusCode)
	})
}
