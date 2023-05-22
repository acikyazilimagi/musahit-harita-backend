package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/cache"
	"github.com/gofiber/fiber/v2"
)

func InvalidateCache() fiber.Handler {
	cacheRepo := cache.NewRedisStore()

	return func(ctx *fiber.Ctx) error {
		err := cacheRepo.DeleteAll()

		if err != nil {
			ctx.Status(fiber.StatusInternalServerError)
			return ctx.SendString(err.Error())
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
