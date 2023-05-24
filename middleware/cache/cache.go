package cache

import (
	"github.com/acikkaynak/musahit-harita-backend/cache"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func New() fiber.Handler {
	cacheRepo, err := cache.NewRedisStore()
	if err != nil {
		log.Logger().Error("cache error", zap.Error(err))
		return nil
	}
	return func(c *fiber.Ctx) error {
		if c.Path() == "/healthcheck" ||
			c.Path() == "/metrics" ||
			c.Path() == "/monitor" ||
			c.Path() == "/swagger/index.html" {
			return c.Next()
		}

		reqURI := c.OriginalURL()
		hashURL := uuid.NewSHA1(uuid.NameSpaceOID, []byte(reqURI)).String()
		if c.Method() != http.MethodGet {
			// Don't cache write endpoints. We can maintain of list to exclude certain http methods later.
			// Since there will be an update in db, better to remove cache entries for this url
			err := cacheRepo.Delete(hashURL)
			if err != nil {
				log.Logger().Error("delete cache error", zap.Error(err))
			}
			return c.Next()
		}
		cacheData := cacheRepo.GetCacheResponse(hashURL)
		if cacheData == nil || len(cacheData) == 0 {
			c.Next()
			if c.Response().StatusCode() == fiber.StatusOK && len(c.Response().Body()) > 0 {
				body, _ := c.Response().BodyUncompressed()
				cacheRepo.SetCacheResponse(hashURL, body, 5*time.Minute)
			}
			return nil
		}

		c.Set("x-cached-response", "true")
		c.Response().SetBodyRaw(cacheData)
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
		return nil
	}
}
